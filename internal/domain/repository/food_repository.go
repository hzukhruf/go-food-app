package repository

import (
	"food-app/internal/domain/entity"

	"github.com/jinzhu/gorm"
)

type FoodRepository struct {
	db *gorm.DB
}

type IFoodRepository interface {
	SaveFood(*entity.Food) (*entity.Food, error)
	GetDetailFood(int) (*entity.Food, error)
	GetAllFood() ([]entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, error)
	DeleteFood(int) error
}

func NewFoodRepository(db *gorm.DB) *FoodRepository {
	var foodRepo = FoodRepository{}
	foodRepo.db = db
	return &foodRepo
}

func (r *FoodRepository) SaveFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Create(&food).Error
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (r *FoodRepository) GetDetailFood(id int) (*entity.Food, error) {
	var food entity.Food
	err := r.db.Where("id = ?", id).Take(&food).Error
	if err != nil {
		return nil, err
	}

	return &food, nil
}

func (r *FoodRepository) GetAllFood() ([]entity.Food, error) {
	var foods []entity.Food
	err := r.db.Order("id desc").Find(&foods).Error
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (r *FoodRepository) UpdateFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Save(&food).Error
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (r *FoodRepository) DeleteFood(id int) error {
	var food entity.Food
	err := r.db.Where("id = ?", id).Delete(&food).Error
	if err != nil {
		return err
	}

	return nil
}
