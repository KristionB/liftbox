package main

import (
	"github.com/gin-gonic/gin"
	"github.com/KristionB/secure-file-sync/server/handlers"
)

func main() {
	r := gin.Default()
	r.POST("/upload", handlers.Upload)
	r.GET("/download/:id", handlers.Download)
	r.Run(":8080")
}

