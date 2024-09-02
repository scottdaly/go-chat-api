package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db = initDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/chat", handleChat)
	r.GET("/history", getHistory)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func handleChat(c *gin.Context) {
	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := callLLMAPI(input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get LLM response"})
		return
	}

	_, err = db.Exec("INSERT INTO chats (user_input, ai_response) VALUES (?, ?)", input.Message, response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func getHistory(c *gin.Context) {
	rows, err := db.Query("SELECT user_input, ai_response, timestamp FROM chats ORDER BY timestamp DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history"})
		return
	}
	defer rows.Close()

	var history []gin.H
	for rows.Next() {
		var userInput, aiResponse string
		var timestamp string
		err := rows.Scan(&userInput, &aiResponse, &timestamp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse chat history"})
			return
		}
		history = append(history, gin.H{
			"user_input":  userInput,
			"ai_response": aiResponse,
			"timestamp":   timestamp,
		})
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
