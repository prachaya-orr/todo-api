package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-api/auth"
	"todo-api/todo"
)

func main() {
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

	r.GET("/tokenz", auth.AccessToken("==signature=="))

	proctected := r.Group("", auth.ProtectMiddleware([]byte("==signature==")))

	todoHandler := todo.NewTodoHandler(db)
	proctected.POST("/todos", todoHandler.NewTask)

	r.Run()
}
