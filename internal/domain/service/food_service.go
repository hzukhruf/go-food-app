package service

import (
	"food-app/internal/domain/entity"
	"food-app/internal/domain/repository"
)


type FoodService struct {
	foodRepo repository.IFoodRepository
	userRepo repository.IUserRepository
}

type IFoodService interface {
	SaveFood(*entity.FoodViewModel) (*entity.FoodViewModel, error)
	GetDetailFood(int) (*entity.FoodDetailViewModel, error)
	GetAllFood() ([]entity.FoodDetailViewModel, error)
	UpdateFood(*entity.FoodViewModel) (*entity.FoodViewModel, error)
	DeleteFood(int) error
}

func NewFoodService(foodRepo repository.IFoodRepository, userRepo repository.IUserRepository) *FoodService {
	var foodService = FoodService{}
	foodService.foodRepo = foodRepo
	foodService.userRepo = userRepo
	return &foodService
}

func (s *FoodService) SaveFood(foodVM *entity.FoodViewModel) (*entity.FoodViewModel, error)  {

	var  food = entity.Food{
		UserID: foodVM.UserID,
		Title: foodVM.Title,
		Description: foodVM.Description,
		FoodImage: foodVM.FoodImage,
	}

	result, err := s.foodRepo.SaveFood(&food)

	if err != nil {
		return nil, err
	}

	if result != nil {
		foodVM = &entity.FoodViewModel{
			ID: result.ID,
			UserID: result.UserID,
			Title: result.Title,
			Description: result.Description,
			FoodImage: result.FoodImage,
		}
	}

	return foodVM, nil
}

func (s *FoodService) GetDetailFood(id int) (*entity.FoodDetailViewModel, error)  {
	result, err := s.foodRepo.GetDetailFood(id)

	if err != nil {
		return nil, err
	}

	var foodVM entity.FoodDetailViewModel

	if result != nil {
		foodVM = entity.FoodDetailViewModel{
			UserName: s.userRepo.GetUserName(result.UserID),
			Title: result.Title,
			FoodImage: result.FoodImage,
			Description: result.Description,
		}
	}

	return &foodVM, nil
}

func (s *FoodService) GetAllFood() ([]entity.FoodDetailViewModel, error)  {
	result, err := s.foodRepo.GetAllFood()
	if err != nil {
		return nil, err
	}

	var foodListVM []entity.FoodDetailViewModel

	if result != nil {
		for _, item := range result {
			foodVM := entity.FoodDetailViewModel{
				UserName: s.userRepo.GetUserName(item.UserID),
				Title: item.Title,
				FoodImage: item.FoodImage,
				Description: item.Description,
			}

			foodListVM = append(foodListVM, foodVM)
		}
	}

	return foodListVM, nil
}

func (s *FoodService) UpdateFood(foodVM *entity.FoodViewModel) (*entity.FoodViewModel, error) {
	var food = entity.Food{
		ID: foodVM.ID,
		UserID: foodVM.UserID,
		Title: foodVM.Title,
		Description: foodVM.Description,
		FoodImage: foodVM.FoodImage,
	}

	_, err := s.foodRepo.UpdateFood(&food)

	if err != nil {
		return nil, err
	}

	return foodVM, nil
}

func (s *FoodService) DeleteFood(id int) error  {
	err := s.foodRepo.DeleteFood(id)
	if err != nil {
		return err
	}

	return nil
}