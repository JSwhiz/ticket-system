package main

import (
    "context"
    "log"
    "time"

	"github.com/joho/godotenv"

    "github.com/JSwhiz/ticket-system/internal/platform/database"
    "github.com/JSwhiz/ticket-system/pkg/config"
    "github.com/gin-gonic/gin"
    "github.com/JSwhiz/ticket-system/internal/app/auth"
    "github.com/JSwhiz/ticket-system/handler"
	ticketHandler "github.com/JSwhiz/ticket-system/handler"
    ticketService "github.com/JSwhiz/ticket-system/internal/app/tickets"

	
)

func main() {

	if err := godotenv.Load(); err != nil {
    log.Printf("No .env file found, relying on environment")
	}

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    db, err := database.Open(cfg.DBURL, "migrations")
    if err != nil {
        log.Fatalf("init database: %v", err)
    }
    defer db.Close()

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

	authSvc := auth.NewService(db, cfg.JWTSecret, cfg.Timeout)

	api := r.Group("/api/v1/auth")
    handler.RegisterAuthRoutes(api, authSvc, cfg)

	protected := r.Group("/api/v1")
    protected.Use(auth.JWTAuthMiddleware(cfg.JWTSecret, db))

	ticketSvc := ticketService.NewService(db)
	ticketHandler.RegisterTicketRoutes(protected, ticketSvc)

    r.GET("/health", func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
        defer cancel()
        if err := db.PingContext(ctx); err != nil {
            c.JSON(500, gin.H{"status": "fail", "error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"status": "ok"})
    })

    //TODO: сюда добавим остальные маршруты /api/v1/…

    addr := ":" + cfg.ServerPort
    log.Printf("Server listening on %s (loglevel=%s)", addr, cfg.LogLevel)
    if err := r.Run(addr); err != nil {
        log.Fatalf("server error: %v", err)
    }

	r.Run(":" + cfg.ServerPort)
}
