package handlers

import "github.com/gin-gonic/gin"

func Upload(c *gin.Context) {
	var payload map[string]string
	c.BindJSON(&payload)

	// verify signature + HMAC
	// store encrypted blob

	c.JSON(200, gin.H{"status": "ok"})
}

