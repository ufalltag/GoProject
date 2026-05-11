package main

import (
	"FirstWorkspace/database"
	"FirstWorkspace/handlers"
	"FirstWorkspace/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	database.RunMigrations()

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

		protected.POST("/analyze", handlers.Analyze)
	}
	r.Run(":8080")
}
