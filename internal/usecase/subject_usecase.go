package usecase

import (
	"context"
	"sample-project/internal/dto"
	"sample-project/internal/entity"
	"sample-project/internal/repository"
)

type SubjectUsecase interface {
	GetSubject(ctx context.Context) ([]entity.Subject, error)
	GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error)
	CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error)
	UpdateSubject(ctx context.Context, id int, updateSubjectDto dto.UpdateSubjectDto) (*entity.Subject, error)
	DeleteSubject(ctx context.Context, id int) error
}

type subjectUseCase struct {
	repo repository.SubjectRepository
}

func NewSubjectUseCase(repo repository.SubjectRepository) SubjectUsecase {
	return &subjectUseCase{repo: repo}
}

func (u *subjectUseCase) GetSubject(ctx context.Context) ([]entity.Subject, error) {
	return u.repo.GetAllSubjects(ctx)
}

func (u *subjectUseCase) GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error) {
	return u.repo.GetSubjectByID(ctx, id)
}

func (u *subjectUseCase) CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error) {
	return u.repo.CreateSubject(ctx, subject)
}

func (u *subjectUseCase) UpdateSubject(ctx context.Context, id int, updateSubjectDto dto.UpdateSubjectDto) (*entity.Subject, error) {
	subject, err := u.repo.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateSubjectDto.Name != "" {
		subject.Name = updateSubjectDto.Name
	}

	if updateSubjectDto.Status {
		subject.Status = updateSubjectDto.Status
	}

	return u.repo.UpdateSubject(ctx, id, *subject)
}

func (u *subjectUseCase) DeleteSubject(ctx context.Context, id int) error {
	return u.repo.DeleteSubject(ctx, id)
}
