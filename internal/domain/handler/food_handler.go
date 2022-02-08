package handler

import (
	"fmt"
	"food-app/internal/domain/entity"
	"food-app/internal/domain/service"
	"food-app/pkg/jwttoken"
	"food-app/pkg/response"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)


type FoodHandler struct {
	foodService service.IFoodService
}

func NewFoodHandler(foodService service.IFoodService) *FoodHandler  {
	var foodHandler = FoodHandler{}
	foodHandler.foodService = foodService
	return &foodHandler
}

func (h *FoodHandler) SaveFood(c *gin.Context)  {
	
	title := c.DefaultPostForm("title","title")
	description := c.DefaultPostForm("description","description")

	tokenMetadata, err := jwttoken.ExtractTokenMetadata(c.Request)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := tokenMetadata.UserID

	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	path := viper.GetString("Files.Path")
	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", path, filename)); err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	var food = entity.FoodViewModel{
		UserID: int(userId),
		Title: title,
		Description: description,
		FoodImage: filename,
	}

	saveFoodError := food.Validate("")
	if len(saveFoodError) > 0 {
		response.ResponseCustomError(c, saveFoodError, http.StatusBadRequest)
		return
	}

	result, err := h.foodService.SaveFood(&food)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
	}

	response.ResponseCreated(c, result)
}

func (h *FoodHandler) GetAllFood(c *gin.Context)  {
	result, err := h.foodService.GetAllFood()
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		result = []entity.FoodDetailViewModel{}
	}

	response.ResponseOKWithData(c, result)
}

func (h *FoodHandler) GetDetailFood(c *gin.Context)  {
	foodId, err := strconv.Atoi(c.Param("food_id"))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.foodService.GetDetailFood(foodId)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResponseOKWithData(c, result)
}

func (h *FoodHandler) UpdateFood(c *gin.Context) {
	foodId, err := strconv.Atoi(c.Param("food_id"))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	title := c.DefaultPostForm("title", "title")
	description := c.DefaultPostForm("description", "description")
	userId := 0

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	var food = entity.FoodViewModel{
		ID:          foodId,
		UserID:      userId,
		Title:       title,
		Description: description,
		FoodImage:   filename,
	}

	saveFoodError := food.Validate("")
	if len(saveFoodError) > 0 {
		response.ResponseCustomError(c, saveFoodError, http.StatusBadRequest)
		return
	}

	result, err := h.foodService.UpdateFood(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FoodHandler) DeleteFood(c *gin.Context) {
	foodId, err := strconv.Atoi(c.Param("food_id"))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.foodService.DeleteFood(foodId)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResponseOK(c, "Successfully Food Deleted")
}