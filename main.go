package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
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

	err = godotenv.Load("./local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s", err)
	}

	dsn := os.Getenv("DB_CONN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	r.Use(cors.New(config))
	// readiness Probe:
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/limitz", limitHandler)
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
	proctected.GET("/todos", todoHandler.List)
	proctected.DELETE("/todos/:id", todoHandler.Remove)

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServeWithGracefulShutdown(s)
}

// maybe use as a middleWare
var limiter = rate.NewLimiter(10, 5)

func limitHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
