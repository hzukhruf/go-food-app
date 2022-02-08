package service

import (
	"fmt"
	"food-app/internal/domain/entity"
	"food-app/internal/domain/repository"
	"food-app/pkg/security"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.IUserRepository
}

type IUserService interface {
	SaveUser(*entity.ReqisterViewModel) (*entity.UserViewModel, error)
	GetListUser() (*[]entity.UserViewModel, error)
	GetDetailUser(id int) (*entity.UserViewModel, error)
	UpdateUser(userVM *entity.User) (*entity.UserViewModel, error)
	DeleteUser(id int) error
	GetUserByEmailPassword(loginVM entity.LoginViewModel) (*entity.User, error)
	SaveUserList(*[]entity.ReqisterViewModel) (*[]entity.UserViewModel, error)
}

func NewUserService(userRepo repository.IUserRepository) *UserService {
	var userService = UserService{}
	userService.userRepo = userRepo
	return &userService
}

func (s *UserService) GetListUser() (*[]entity.UserViewModel, error) {
	result, err := s.userRepo.GetAllUser()
	if err != nil {
		return nil, err
	}

	var users []entity.UserViewModel
	for _, item := range result {
		var user entity.UserViewModel
		user.Email = item.Email
		user.FullName = fmt.Sprintf("%s %s", item.FirstName, item.LastName)
		user.Email = item.Email
		users = append(users, user)
	}

	return &users, nil
}

func (s *UserService) GetDetailUser(id int) (*entity.UserViewModel, error) {
	var viewModel entity.UserViewModel

	result, err := s.userRepo.GetDetailUser(id)
	if err != nil {
		return nil, err
	}

	if result != nil {
		viewModel = entity.UserViewModel{
			ID:       result.ID,
			FullName: fmt.Sprintf("%s %s", result.FirstName, result.LastName),
			Email:    result.Email,
		}
	}

	return &viewModel, nil
}

func (s *UserService) UpdateUser(userVM *entity.User) (*entity.UserViewModel, error) {
	password, err := userVM.EncryptPassword(userVM.Password)
	if err != nil {
		return nil, err
	}

	userVM.Password = password

	result, err := s.userRepo.UpdateUser(userVM)
	if err != nil {
		return nil, err
	}

	var userAfterUpdate entity.UserViewModel
	userAfterUpdate = entity.UserViewModel{
		ID:       result.ID,
		FullName: fmt.Sprintf("%s %s", result.FirstName, result.LastName),
		Email:    result.Email,
	}

	return &userAfterUpdate, err
}

func (s *UserService) DeleteUser(id int) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserByEmailPassword(loginVM entity.LoginViewModel) (*entity.User, error) {
	result, err := s.userRepo.GetUserByEmailPassword(loginVM)
	if err != nil {
		return nil, err
	}

	// Verify Password
	err = security.VerifyPassword(result.Password, loginVM.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("Incorrect Password. Error %s", err.Error())
	}

	return result, nil
}

func (s *UserService) SaveUser(userVM *entity.ReqisterViewModel) (*entity.UserViewModel, error) {
	var user = entity.User{
		FirstName: userVM.FirstName,
		LastName:  userVM.LastName,
		Email:     userVM.Email,
	}

	password, err := user.EncryptPassword(userVM.Password)
	if err != nil {
		return nil, err
	}

	user.Password = password

	result, err := s.userRepo.SaveUser(&user)
	if err != nil {
		return nil, err
	}

	var afterRegVM entity.UserViewModel

	if result != nil {
		afterRegVM = entity.UserViewModel{
			ID:       result.ID,
			FullName: fmt.Sprintf("%s %s", result.FirstName, result.LastName),
			Email:    result.Email,
		}
	}

	return &afterRegVM, nil
}

/* untuk contoh go routine */
func (s *UserService) SaveUserList(userVmList *[]entity.ReqisterViewModel) (*[]entity.UserViewModel, error) {

	n := len(*userVmList) / 2

	var chann = make(chan *entity.UserViewModel)

	go addUser(0, n, *userVmList, chann, s)

	go addUser(n, len(*userVmList), *userVmList, chann, s)

	var userViewModelList []entity.UserViewModel
	for i := 0; i < len(*userVmList); i++ {
		var userVm = <-chann
		newUserVM := entity.UserViewModel{
			ID:       userVm.ID,
			FullName: userVm.FullName,
			Email:    userVm.Email,
		}
		userViewModelList = append(userViewModelList, newUserVM)
	}

	return &userViewModelList, nil
}

func addUser(awal int, akhir int, userVmList []entity.ReqisterViewModel, chann chan *entity.UserViewModel, s *UserService) {
	for i := awal; i < akhir; i++ {
		userView, _ := s.SaveUser(&userVmList[i])
		chann <- userView
	}
}

/* untuk contoh go routine */
