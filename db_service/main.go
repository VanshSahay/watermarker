package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
    ID              uint   `gorm:"primaryKey"`
    Title           string `json:"title"`
    Author          string `json:"author"`
    WatermarkStatus string `json:"watermark_status"`
    Watermark       string `json:"watermark"`
}

var db *gorm.DB

func main() {
    var err error
    dsn := "host=localhost user=postgres password=yourpassword dbname=watermarker port=5432 sslmode=disable"
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    db.AutoMigrate(&Document{})

    r := gin.Default()

    r.GET("/documents", func(c *gin.Context) {
        var docs []Document
        db.Find(&docs)
        c.JSON(http.StatusOK, docs)
    })

    r.POST("/documents", func(c *gin.Context) {
        var doc Document
        if err := c.BindJSON(&doc); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }
        db.Create(&doc)
        c.JSON(http.StatusCreated, doc)
    })

    r.Run(":8082")
}
