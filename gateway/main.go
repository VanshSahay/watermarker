package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func proxyRequest(c *gin.Context, targetURL string) {
	client := &http.Client{}

	// Log the incoming request body
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for reuse
	c.Writer.WriteString("Forwarding request to " + targetURL + "\n")
	c.Writer.WriteString("Request Body: " + string(body) + "\n")

	// Create the forwarded request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header = c.Request.Header

	// Log the headers
	for key, values := range req.Header {
		for _, value := range values {
			c.Writer.WriteString("Header: " + key + " -> " + value + "\n")
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read and log the response
	respBody, _ := io.ReadAll(resp.Body)
	c.Writer.WriteString("Response Body: " + string(respBody) + "\n")
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}


func main() {
	r := gin.Default()

	// Route for the Authorization Microservice
	r.POST("/auth", func(c *gin.Context) {
		proxyRequest(c, "http://localhost:8081/auth")
	})

	// Route for the Database Microservice
	r.GET("/documents", func(c *gin.Context) {
		proxyRequest(c, "http://localhost:8082/documents")
	})
	r.POST("/documents", func(c *gin.Context) {
		proxyRequest(c, "http://localhost:8082/documents")
	})

	// Route for the Watermark Microservice
	r.POST("/watermark", func(c *gin.Context) {
		proxyRequest(c, "http://localhost:8083/watermark")
	})

	r.Run(":8080") // Run the gateway on port 8080
}
