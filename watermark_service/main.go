package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WatermarkRequest struct {
	DocumentID string `json:"document_id" binding:"required"`
	Watermark  string `json:"watermark" binding:"required"`  
}


func main() {
    r := gin.Default()

    r.POST("/watermark", func(c *gin.Context) {
        // Simulate watermarking logic
        var req WatermarkRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            fmt.Println(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "Watermarking started", "document": req})
    })

    r.Run(":8083")
}
