package main

import (
	"FirstWorkspace/database"
	"FirstWorkspace/handlers"
	"FirstWorkspace/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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
		protected.PUT("/change-password", handlers.ChangePassword)

		protected.GET("/folders", handlers.GetFolders)
		protected.POST("/folders", handlers.CreateFolder)
		protected.DELETE("/folders/:id", handlers.DeleteFolder)
		protected.PUT("/folders/:id", handlers.UpdateFolder)

		protected.GET("/folders/:id/bookmarks", handlers.GetBookmarks)
		protected.GET("/bookmarks/recent", handlers.GetRecentBookmarks)
		protected.POST("/bookmarks", handlers.CreateBookmark)
		protected.DELETE("/bookmarks/:id", handlers.DeleteBookmark)
		protected.PUT("/bookmarks/:id", handlers.UpdateBookmark)

		protected.POST("/analyze", handlers.Analyze)
	}
	r.Run(":8080")
}
