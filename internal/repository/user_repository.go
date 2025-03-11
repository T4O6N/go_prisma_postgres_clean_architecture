package repository

import (
	"context"
	"fmt"
	"sample-project/internal/entity"
	"sample-project/internal/utils"
	"sample-project/prisma/db"
	"time"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	client *db.PrismaClient
}

func NewUserRepository(client *db.PrismaClient) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users, err := r.client.User.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []entity.User
	for _, u := range users {
		var subjectID int
		if id, ok := u.SubjectID(); ok {
			subjectID = id
		}

		result = append(result, entity.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			SubjectID: subjectID,
			Status:    u.Status,
			CreatedAt: utils.FormatToVientianeTime(u.CreatedAt),
			UpdatedAt: utils.FormatToVientianeTime(u.UpdatedAt),
		})
	}
	return result, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: utils.FormatToVientianeTime(user.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(user.UpdatedAt),
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// Check if subject_id exists if it's provided and not zero
	if user.SubjectID != 0 {
		// Check if the subject exists
		subject, err := r.client.Subject.FindUnique(
			db.Subject.ID.Equals(user.SubjectID),
		).Exec(ctx)

		if err != nil || subject == nil {
			return nil, fmt.Errorf("subject with ID: %d not found", user.SubjectID)
		}
	}

	newUser, err := r.client.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.SubjectID.Set(user.SubjectID),
		db.User.Status.Set(user.Status),
		db.User.CreatedAt.Set(utils.FormatToVientianeTime(time.Now())),
		db.User.UpdatedAt.Set(utils.FormatToVientianeTime(time.Now())),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		SubjectID: user.SubjectID,
		CreatedAt: utils.FormatToVientianeTime(newUser.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(newUser.UpdatedAt),
	}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error) {
	if user.SubjectID != 0 {
		subject, err := r.client.Subject.FindUnique(
			db.Subject.ID.Equals(user.SubjectID),
		).Exec(ctx)

		if err != nil || subject == nil {
			return nil, fmt.Errorf("subject with ID: %d not found", user.SubjectID)
		}
	}

	updateUser, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.SubjectID.Set(user.SubjectID),
		db.User.Status.Set(user.Status),
		db.User.UpdatedAt.Set(utils.FormatToVientianeTime(time.Now())),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        updateUser.ID,
		Name:      updateUser.Name,
		Email:     updateUser.Email,
		SubjectID: user.SubjectID,
		Status:    updateUser.Status,
		CreatedAt: utils.FormatToVientianeTime(updateUser.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(updateUser.UpdatedAt),
	}, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(ctx)

	return err
}
