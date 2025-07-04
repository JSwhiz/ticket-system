package main

import (
    "context"
    "log"
    "time"

    "github.com/joho/godotenv"
    "github.com/JSwhiz/ticket-system/handler"
    "github.com/JSwhiz/ticket-system/internal/app/auth"
    "github.com/JSwhiz/ticket-system/internal/app/authz"
    "github.com/JSwhiz/ticket-system/internal/app/refs"
    "github.com/JSwhiz/ticket-system/internal/app/tickets"
    "github.com/JSwhiz/ticket-system/internal/platform/database"
    "github.com/JSwhiz/ticket-system/pkg/config"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func main() {
    // 1) Загрузка .env (опционально)
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found, relying on env vars")
    }

    // 2) Конфиг
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    // 3) Открываем БД и прогоняем миграции
    db, err := database.Open(cfg.DBURL, "migrations")
    if err != nil {
        log.Fatalf("init database: %v", err)
    }
    defer db.Close()

    // 4) Инициализируем сервисы
    authSvc := auth.NewService(db, cfg.JWTSecret, cfg.Timeout)  // <- cfg.TokenTTL
    authzRepo := authz.NewRepository(db)
    ticketSvc := tickets.NewService(db)
    refsRepo  := refs.NewRepository(db)
    refsSvc   := refs.NewService(refsRepo)

    // 5) Настраиваем Gin
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    // 6) Public — только логин
    authGroup := r.Group("/api/v1/auth")
    handler.RegisterAuthRoutes(authGroup, authSvc, cfg)

    // 7) Protected — все остальные роуты
    protected := r.Group("/api/v1")
    protected.Use(auth.JWTAuthMiddleware(cfg.JWTSecret))

    //   7.1 Справочники (GET-only)
    handler.RegisterRefRoutes(protected, refsSvc)
	handler.RegisterTicketRoutes(protected, ticketSvc, authzRepo)

    // 8) Healthcheck
    r.GET("/health", func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
        defer cancel()
        if err := db.PingContext(ctx); err != nil {
            c.JSON(500, gin.H{"status": "fail", "error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"status": "ok"})
    })

    // 9) Стартуем сервер
    addr := ":" + cfg.ServerPort
    log.Printf("Server listening on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
