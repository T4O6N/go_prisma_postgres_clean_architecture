package usecase

import (
	"context"
	"sample-project/internal/entity"
	"sample-project/internal/repository"
)

// NOTE - subject use case interface
type SubjectUsecase interface {
	GetSubject(ctx context.Context) ([]entity.Subject, error)
	GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error)
	CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error)
	UpdateSubject(ctx context.Context, id int, updateSubject entity.Subject) (*entity.Subject, error)
	DeleteSubject(ctx context.Context, id int) error
	ClearSubjectCache(ctx context.Context) error
}

// NOTE - subject use case struct
type subjectUseCase struct {
	repo repository.SubjectRepository
}

// NOTE - new subject use case
func NewSubjectUseCase(repo repository.SubjectRepository) SubjectUsecase {
	return &subjectUseCase{repo: repo}
}

// NOTE - get all subjects use case
func (u *subjectUseCase) GetSubject(ctx context.Context) ([]entity.Subject, error) {
	return u.repo.GetAllSubjects(ctx)
}

// NOTE - get subject by id use case
func (u *subjectUseCase) GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error) {
	return u.repo.GetSubjectByID(ctx, id)
}

// NOTE - create subject use case
func (u *subjectUseCase) CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error) {
	return u.repo.CreateSubject(ctx, subject)
}

// NOTE - update subject use case
func (u *subjectUseCase) UpdateSubject(ctx context.Context, id int, updateSubject entity.Subject) (*entity.Subject, error) {
	subject, err := u.repo.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateSubject.Name != "" {
		subject.Name = updateSubject.Name
	}

	subject.Status = updateSubject.Status

	return u.repo.UpdateSubject(ctx, id, *subject)
}

// NOTE - delete subject use case
func (u *subjectUseCase) DeleteSubject(ctx context.Context, id int) error {
	return u.repo.DeleteSubject(ctx, id)
}

// NOTE - clear subjects cache use case
func (u *subjectUseCase) ClearSubjectCache(ctx context.Context) error {
	return u.repo.ClearSubjectCache(ctx)
}
