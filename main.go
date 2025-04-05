package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-api/auth"
	"todo-api/server"
	"todo-api/todo"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	// liveness Probe:  Cat file in kub better than http check. In case, service get SIGTERM
	// if our service not receive request from client because service prepare for graceful shut down
	//
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("please consider environment variables: %s", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	// readiness Probe:
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buidcommit": buildcommit,
			"buildtime":  buildtime,
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))
	proctected := r.Group("", auth.ProtectMiddleware([]byte(os.Getenv("SIGN"))))

	todoHandler := todo.NewTodoHandler(db)
	proctected.POST("/todos", todoHandler.NewTask)

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServeWithGracefulShutdown(s)
}
