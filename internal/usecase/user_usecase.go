package usecase

import (
	"context"
	"sample-project/internal/entity"
	"sample-project/internal/repository"
)

// NOTE - user use case interface
type UserUseCase interface {
	GetUsers(ctx context.Context, page, limit int, name string, startDate, endDate string) ([]entity.User, int, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	ClearUserCache(ctx context.Context) error
}

// NOTE - user use case struct
type userUsecase struct {
	repo repository.UserRepository
}

// NOTE - new user use case
func NewUserUsecase(repo repository.UserRepository) UserUseCase {
	return &userUsecase{repo: repo}
}

// NOTE - get all users use case
func (u *userUsecase) GetUsers(ctx context.Context, page, limit int, name string, startDate, endDate string) ([]entity.User, int, error) {
	return u.repo.GetAllUsers(ctx, page, limit, name, startDate, endDate)
}

// NOTE - get user by id use case
func (u *userUsecase) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

// NOTE - create user use case
func (u *userUsecase) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	return u.repo.CreateUser(ctx, user)
}

// NOTE - update user use case
func (u *userUsecase) UpdateUser(ctx context.Context, id int, userUpdate entity.User) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}
	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}

	user.Status = userUpdate.Status

	if userUpdate.SubjectID != 0 {
		user.SubjectID = userUpdate.SubjectID
	}

	return u.repo.UpdateUser(ctx, id, *user)
}

// NOTE - delete user use case
func (u *userUsecase) DeleteUser(ctx context.Context, id int) error {
	return u.repo.DeleteUser(ctx, id)
}

// NOTE - clear users cache use case
func (u *userUsecase) ClearUserCache(ctx context.Context) error {
	return u.repo.ClearUserCache(ctx)
}
