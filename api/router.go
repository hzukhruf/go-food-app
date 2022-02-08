package api

import (
	"food-app/api/middleware"
	"food-app/internal/domain/handler"
	"food-app/internal/domain/repository"
	"food-app/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	router := gin.Default()

	//Register User Repo
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//Register Food Repo
	foodRepo := repository.NewFoodRepository(db)
	foodService := service.NewFoodService(foodRepo, userRepo)
	foodHandler := handler.NewFoodHandler(foodService)

	food := router.Group("v1/food")
	{
		food.GET("/", middleware.AuthMiddleware(), middleware.CORSMiddleware(), foodHandler.GetAllFood)
		food.GET("/:food_id", middleware.AuthMiddleware(), middleware.CORSMiddleware(), foodHandler.GetDetailFood)
		food.POST("/", middleware.AuthMiddleware(), middleware.CORSMiddleware(), foodHandler.SaveFood)
		food.PUT("/:food_id", middleware.CORSMiddleware(), foodHandler.UpdateFood)
		food.DELETE("/:food_id", middleware.CORSMiddleware(), foodHandler.DeleteFood)
	}

	user := router.Group("v1/user")
	{
		user.GET("/", middleware.CORSMiddleware(), userHandler.GetAllUser)
		user.GET("/:user_id", middleware.AuthMiddleware(), middleware.CORSMiddleware(), userHandler.GetDetailUser)
		user.POST("/", middleware.CORSMiddleware(), userHandler.RegisterUser)
		user.PUT("/:user_id", middleware.CORSMiddleware(), userHandler.UpdateUser)
		user.DELETE("/:user_id", middleware.CORSMiddleware(), userHandler.DeleteUser)
		user.POST("/login", middleware.CORSMiddleware(), userHandler.Login)

		/* untuk contoh go routine */
		user.POST("/userList", middleware.CORSMiddleware(), userHandler.RegisterUserList)
	}

	return router

}
