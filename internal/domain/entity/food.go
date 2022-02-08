package entity

import "strings"

type Food struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FoodImage   string `json:"food_image"`
}

type FoodDetailViewModel struct {
	Title       string `json"title"`
	Description string `json:"description"`
	FoodImage   string `json:"food_image"`
	UserName    string `json:"user_name"`
}

type FoodViewModel struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json"title"`
	Description string `json:"description"`
	FoodImage   string `json:"food_image"`
}

func (f *FoodViewModel) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	default:
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}