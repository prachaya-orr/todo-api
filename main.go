package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-api/auth"
	"todo-api/todo"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("please consider environment variables: %s", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	proctected := r.Group("", auth.ProtectMiddleware([]byte(os.Getenv("SIGN"))))

	todoHandler := todo.NewTodoHandler(db)
	proctected.POST("/todos", todoHandler.NewTask)

	r.Run()
}
