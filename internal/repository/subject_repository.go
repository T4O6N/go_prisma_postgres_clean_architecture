package repository

import (
	"context"
	"sample-project/internal/entity"
	"sample-project/internal/utils"
	"sample-project/prisma/db"
	"time"
)

type SubjectRepository interface {
	GetAllSubjects(ctx context.Context) ([]entity.Subject, error)
	GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error)
	CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error)
	UpdateSubject(ctx context.Context, id int, subject entity.Subject) (*entity.Subject, error)
	DeleteSubject(ctx context.Context, id int) error
}

type subjectRepository struct {
	client *db.PrismaClient
}

func NewSubjectRepository(client *db.PrismaClient) SubjectRepository {
	return &subjectRepository{client: client}
}

func (r *subjectRepository) GetAllSubjects(ctx context.Context) ([]entity.Subject, error) {
	subjects, err := r.client.Subject.FindMany().With(
		db.Subject.User.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []entity.Subject
	for _, s := range subjects {
		var users []entity.User
		for _, u := range s.User() {
			users = append(users, entity.User{
				ID:        u.ID,
				Name:      u.Name,
				Email:     u.Email,
				Status:    u.Status,
				CreatedAt: utils.FormatToVientianeTime(u.CreatedAt),
				UpdatedAt: utils.FormatToVientianeTime(u.UpdatedAt),
			})
		}
		result = append(result, entity.Subject{
			ID:        s.ID,
			Name:      s.Name,
			User:      users,
			Status:    s.Status,
			CreatedAt: utils.FormatToVientianeTime(s.CreatedAt),
			UpdatedAt: utils.FormatToVientianeTime(s.UpdatedAt),
		})
	}
	return result, nil
}

func (r *subjectRepository) GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error) {
	subject, err := r.client.Subject.FindUnique(
		db.Subject.ID.Equals(id),
	).With(
		db.Subject.User.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for _, u := range subject.User() {
		users = append(users, entity.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: utils.FormatToVientianeTime(u.CreatedAt),
			UpdatedAt: utils.FormatToVientianeTime(u.UpdatedAt),
		})
	}

	return &entity.Subject{
		ID:        subject.ID,
		Name:      subject.Name,
		User:      users,
		Status:    subject.Status,
		CreatedAt: utils.FormatToVientianeTime(subject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(subject.UpdatedAt),
	}, nil
}

func (r *subjectRepository) CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error) {
	newSubject, err := r.client.Subject.CreateOne(
		db.Subject.Name.Set(subject.Name),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &entity.Subject{
		ID:        newSubject.ID,
		Name:      newSubject.Name,
		Status:    newSubject.Status,
		CreatedAt: utils.FormatToVientianeTime(newSubject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(newSubject.UpdatedAt),
	}, nil
}

func (r *subjectRepository) UpdateSubject(ctx context.Context, id int, subject entity.Subject) (*entity.Subject, error) {
	updateSubject, err := r.client.Subject.FindUnique(
		db.Subject.ID.Equals(id),
	).Update(
		db.Subject.Name.Set(subject.Name),
		db.Subject.Status.Set(subject.Status),
		db.Subject.UpdatedAt.Set(utils.FormatToVientianeTime(time.Now())),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &entity.Subject{
		ID:        updateSubject.ID,
		Name:      updateSubject.Name,
		Status:    updateSubject.Status,
		CreatedAt: utils.FormatToVientianeTime(updateSubject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(updateSubject.UpdatedAt),
	}, nil
}

func (r *subjectRepository) DeleteSubject(ctx context.Context, id int) error {
	_, err := r.client.Subject.FindUnique(
		db.Subject.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}
