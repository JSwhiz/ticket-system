package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JSwhiz/ticket-system/internal/app/auth"
	"github.com/JSwhiz/ticket-system/internal/app/tickets"
	"github.com/gin-gonic/gin"
)

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
    if createdFrom := c.Query("created_from"); createdFrom != "" {
        if t, err := time.Parse(time.RFC3339, createdFrom); err == nil {
            filters["created_from"] = t
        }
    }
    if createdTo := c.Query("created_to"); createdTo != "" {
        if t, err := time.Parse(time.RFC3339, createdTo); err == nil {
            filters["created_to"] = t
        }
    }
    if search := c.Query("search"); search != "" {
        filters["search"] = search
    }

    // для пагинации (если будет использоваться)
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

func RegisterTicketRoutes(rg *gin.RouterGroup, svc *tickets.Service) {
    t := rg.Group("/tickets")

    t.GET("", func(c *gin.Context) {
        params := parseListParams(c)
        out, err := svc.List(params)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, out)
    })

    t.POST("", func(c *gin.Context) {
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
    })

    t.GET("/:id", func(c *gin.Context) {
        tkt, err := svc.Get(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, tkt)
    })

    t.PUT("/:id", func(c *gin.Context) {
        var ut tickets.UpdateTicket
        if err := c.ShouldBindJSON(&ut); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        tkt, err := svc.Update(c.Param("id"), ut)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, tkt)
    })

    t.DELETE("/:id", func(c *gin.Context) {
        if err := svc.Delete(c.Param("id")); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    })
}

