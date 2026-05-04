package main

import (
	"FirstWorkspace/database"
	"FirstWorkspace/handlers"
	"FirstWorkspace/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	dsn := fmt.Sprintf(
		"%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Запуск миграции
	database.RunMigrations(dsn)

	database.Connect()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/refresh", handlers.RefreshToken)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired)
	{
		protected.GET("/profile", handlers.Profile)

		protected.GET("/folders", handlers.GetFolders)
		protected.POST("/folders", handlers.CreateFolder)
		protected.DELETE("/folders/:id", handlers.DeleteFolder)
		protected.PUT("/folders/:id", handlers.UpdateFolder)

		protected.GET("/folders/:id/bookmarks", handlers.GetBookmarks)
		protected.POST("/bookmarks", handlers.CreateBookmark)
		protected.DELETE("/bookmarks/:id", handlers.DeleteBookmark)
		protected.PUT("/bookmarks/:id", handlers.UpdateBookmark)
	}
	r.Run(":8080")
}
