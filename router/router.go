package router

import (
	"expvar"
	"gins/controller"
	"gins/middleware"
	"gins/repository"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userRepository repository.UsersRepository,
	tagsController *controller.TagsController,
	authenticationController *controller.AuthenticationController,
	usersController *controller.UserController) *gin.Engine {
	service := gin.Default()

	//  logging and recovery middleware
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(middleware.RateLimitMiddleware(1, 5))
	service.Use(middleware.MetricsMiddleware())

	// CORS middleware
	service.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	service.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Welcome Home!")
	})

	// API group
	router := service.Group("/api")

	// Authentication routes
	authenticationRouter := router.Group("/auth")
	authenticationRouter.POST("/register", authenticationController.Register)
	authenticationRouter.POST("/login", authenticationController.Login)

	// Protected routes: require user to be authenticated
	router.Use(middleware.DeserializeUser(userRepository))

	// Tags routes
	tagsRouter := router.Group("/tags")
	{
		tagsRouter.GET("", tagsController.FindAll)
		tagsRouter.GET("/:tagId", tagsController.FindById)
		tagsRouter.POST("", tagsController.Create)
		tagsRouter.PUT("/:tagId", tagsController.Update)
		tagsRouter.DELETE("/:tagId", tagsController.Delete)
	}

	// Users routes
	usersRouter := router.Group("/users")
	{
		usersRouter.GET("", usersController.GetUsers)
	}

	// Metrics
	service.GET("/debug/vars", gin.WrapH(expvar.Handler()))

	return service
}
