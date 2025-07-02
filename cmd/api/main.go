package main

import (
	"fmt"
	"ticket-system/internal/config"
	"ticket-system/internal/db"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"
	"ticket-system/internal/repository"
	"ticket-system/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	// Подключаемся к базе данных
	dbConn, err := db.Connect(cfg.DBURL)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer dbConn.Close()

	fmt.Println("Успешно подключились к базе данных!")

	// Настраиваем репозитории, сервисы и хендлеры
	userRepo := repository.NewUserRepository(dbConn)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	ticketRepo := repository.NewTicketRepository(dbConn)
	ticketService := services.NewTicketService(ticketRepo)
	ticketHandler := handlers.NewTicketHandler(ticketService, authService)

	// Настраиваем HTTP-сервер
	router := gin.Default()

	// Регистрируем маршруты
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Маршруты для авторизации
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Защищенные маршруты для тикетов
	tickets := router.Group("/api/v1/tickets")
	{
		tickets.Use(middleware.AuthMiddleware(cfg))
		tickets.POST("", ticketHandler.CreateTicket)
		tickets.GET("/:id", ticketHandler.GetTicket)
		tickets.GET("", ticketHandler.GetAllTickets)
		tickets.PATCH("/:id", ticketHandler.UpdateTicket)
		tickets.DELETE("/:id", ticketHandler.DeleteTicket)
	}

	// Запускаем сервер
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
