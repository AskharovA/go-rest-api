package routes

import (
	"AskharovA/go-rest-api/middlewares"
	"AskharovA/go-rest-api/repositories"
	"AskharovA/go-rest-api/services"
	"database/sql"

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
	authenticated.POST("/events/:id/register", func(c *gin.Context) { registerForEvent(c, dbConn) })
	authenticated.DELETE("/events/:id/register", func(c *gin.Context) { cancelRegistration(c, dbConn) })

	server.POST("/signup", userAPI.signup)
	server.POST("/login", userAPI.login)
}
