package main

import (
	"fmt"
	"log"

	"mysql-ip-connect/golang/internal/config"
	"mysql-ip-connect/golang/internal/db"
	"mysql-ip-connect/golang/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Init(cfg.MySQLUsername, cfg.MySQLPassword, cfg.MySQLAddress, cfg.MySQLDatabase, cfg.TableName)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	s := handlers.NewServer(database)

	r.GET("/", s.GetIndex)
	r.GET("/api/status", s.GetStatus) // 新增状态接口
	r.GET("/api/count", s.GetCount)
	r.POST("/api/count", s.PostCount)
	r.GET("/api/wx_openid", s.WxOpenID)

	port := ":" + cfg.Port
	fmt.Printf("Starting server on port %s\n", port)

	if database.IsConnected() {
		fmt.Printf("✅ Database connected successfully\n")
	} else {
		fmt.Printf("⚠️  Warning: Database connection failed, application will run with limited functionality\n")
		fmt.Printf("Please check your .env file configuration\n")
	}

	if err := r.Run(port); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
