package handler

import (
    "io"
    "net/http"
    "strconv"
    "time"

    "github.com/JSwhiz/ticket-system/internal/app/auth"
    "github.com/JSwhiz/ticket-system/internal/app/authz"
    "github.com/JSwhiz/ticket-system/internal/app/tickets"
    "github.com/gin-gonic/gin"
)

// parseListParams разбирает query-параметры для /tickets
func parseListParams(c *gin.Context) map[string]interface{} {
    filters := make(map[string]interface{})

    if s := c.Query("status_id"); s != "" {
        if id, err := strconv.Atoi(s); err == nil {
            filters["status_id"] = id
        }
    }
    if s := c.Query("priority_id"); s != "" {
        if id, err := strconv.Atoi(s); err == nil {
            filters["priority_id"] = id
        }
    }
    if a := c.Query("assignee_id"); a != "" {
        filters["assignee_id"] = a
    }
    if d := c.Query("department_id"); d != "" {
        if id, err := strconv.Atoi(d); err == nil {
            filters["department_id"] = id
        }
    }
    if cf := c.Query("created_from"); cf != "" {
        if t, err := time.Parse(time.RFC3339, cf); err == nil {
            filters["created_from"] = t
        }
    }
    if ct := c.Query("created_to"); ct != "" {
        if t, err := time.Parse(time.RFC3339, ct); err == nil {
            filters["created_to"] = t
        }
    }
    if search := c.Query("search"); search != "" {
        filters["search"] = search
    }
    if p := c.Query("page"); p != "" {
        if page, err := strconv.Atoi(p); err == nil {
            filters["page"] = page
        }
    }
    if ps := c.Query("page_size"); ps != "" {
        if size, err := strconv.Atoi(ps); err == nil {
            filters["page_size"] = size
        }
    }

    return filters
}

func RegisterTicketRoutes(
    rg *gin.RouterGroup,
    svc *tickets.Service,
    authzRepo *authz.Repository,
) {
    t := rg.Group("/tickets")

    // List
    t.GET("",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            params := parseListParams(c)
            out, err := svc.List(params)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, out)
        },
    )

    // Create
    t.POST("",
        authz.Authorize(authzRepo, "create_ticket"),
        func(c *gin.Context) {
            var nt tickets.NewTicket
            if err := c.ShouldBindJSON(&nt); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
            creatorID := c.GetString(auth.CtxUserIDKey)
            tkt, err := svc.Create(nt, creatorID)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusCreated, tkt)
        },
    )

    // Get by ID
    t.GET("/:id",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            tkt, err := svc.Get(c.Param("id"))
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
                return
            }
            c.JSON(http.StatusOK, tkt)
        },
    )

    // Update
    t.PUT("/:id",
        authz.Authorize(authzRepo, "update_ticket"),
        func(c *gin.Context) {
            var in tickets.UpdateTicket
            if err := c.ShouldBindJSON(&in); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
            // минимум одно поле для обновления
            if in.Title == nil &&
                in.Description == nil &&
                in.StatusID == nil &&
                in.PriorityID == nil &&
                in.AssigneeID == nil &&
                in.DepartmentID == nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "at least one field must be provided to update",
                })
                return
            }
            userID := c.GetString(auth.CtxUserIDKey)
            tkt, err := svc.Update(c.Param("id"), in, userID)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, tkt)
        },
    )

    // Delete (soft)
    t.DELETE("/:id",
        authz.Authorize(authzRepo, "delete_ticket"),
        func(c *gin.Context) {
            if err := svc.Delete(c.Param("id")); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.Status(http.StatusNoContent)
        },
    )

    // History
    t.GET("/:id/history",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            hist, err := svc.History(c.Param("id"))
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, hist)
        },
    )

    // Comments
    t.GET("/:id/comments",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            comms, err := svc.Comments(c.Param("id"))
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, comms)
        },
    )
    t.POST("/:id/comments",
        authz.Authorize(authzRepo, "update_ticket"),
        func(c *gin.Context) {
            var in tickets.NewComment
            if err := c.ShouldBindJSON(&in); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
            authorID := c.GetString(auth.CtxUserIDKey)
            cm, err := svc.AddComment(c.Param("id"), authorID, in.Content)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusCreated, cm)
        },
    )

    // Attachments
    t.GET("/:id/attachments",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            atts, err := svc.Attachments(c.Param("id"))
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusOK, atts)
        },
    )
    t.POST("/:id/attachments",
        authz.Authorize(authzRepo, "update_ticket"),
        func(c *gin.Context) {
            file, err := c.FormFile("file")
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
                return
            }
            f, err := file.Open()
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            defer f.Close()
            data, err := io.ReadAll(f)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            authorID := c.GetString(auth.CtxUserIDKey)
            att, err := svc.UploadAttachment(c.Param("id"), authorID, file.Filename, data)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            c.JSON(http.StatusCreated, att)
        },
    )
    t.GET("/:id/attachments/:aid",
        authz.Authorize(authzRepo, "view_ticket"),
        func(c *gin.Context) {
            filename, data, err := svc.DownloadAttachment(c.Param("aid"))
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
                return
            }
            c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
            c.Data(http.StatusOK, "application/octet-stream", data)
        },
    )
}
