package main

import (
	"AskharovA/go-rest-api/config"
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/routes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(dbConn *sql.DB) *gin.Engine {
	router := gin.Default()
	routes.RegisterRoutes(router, dbConn)

	return router
}

func main() {
	dbConn, err := db.InitDB("api.db")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	defer dbConn.Close()

	err = db.CreateTables(dbConn)
	if err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	router := setupRouter(dbConn)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().AppPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %d\n", config.Get().AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
