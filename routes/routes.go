package routes

import (
	"AskharovA/go-rest-api/middlewares"
	"AskharovA/go-rest-api/repositories"
	"AskharovA/go-rest-api/services"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, dbConn *sql.DB) {
	eventRepo := repositories.NewEventRepository(dbConn)
	eventService := services.NewEventService(eventRepo)
	eventAPI := &EventAPI{EventService: eventService}

	userRepo := repositories.NewUsersRepository(dbConn)
	userService := services.NewUserService(userRepo)
	userAPI := &UserAPI{UserService: *userService}

	server.GET("/events", eventAPI.getEvents)
	server.GET("/events/:id", eventAPI.getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", eventAPI.createEvent)
	authenticated.PUT("/events/:id", eventAPI.updateEvent)
	authenticated.DELETE("/events/:id", eventAPI.deleteEvent)
	authenticated.POST("/events/:id/register", eventAPI.registerForEvent)
	authenticated.DELETE("/events/:id/register", eventAPI.cancelRegistration)

	server.POST("/signup", userAPI.signup)
	server.POST("/login", userAPI.login)

	// CORS Middleware Config
	// Allow all origins
	server.Use(cors.Default())

	// Detailed CORS config
	// server.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"https://foo.com", "http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
}
