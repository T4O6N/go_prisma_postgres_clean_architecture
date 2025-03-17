package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sample-project/internal/config/cache"
	"sample-project/internal/entity"
	"sample-project/internal/utils"
	"sample-project/prisma/db"
	"time"

	"github.com/redis/go-redis/v9"
)

// NOTE - subject repository interface
type SubjectRepository interface {
	GetAllSubjects(ctx context.Context) ([]entity.Subject, error)
	GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error)
	CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error)
	UpdateSubject(ctx context.Context, id int, subject entity.Subject) (*entity.Subject, error)
	DeleteSubject(ctx context.Context, id int) error
	ClearSubjectCache(ctx context.Context) error
}

// NOTE - subject repository struct
type subjectRepository struct {
	client      *db.PrismaClient
	redisClient *redis.Client
}

func NewSubjectRepository(client *db.PrismaClient, redisClient *redis.Client) SubjectRepository {
	return &subjectRepository{client: client, redisClient: redisClient}
}

// NOTE - get all subjects repository
func (r *subjectRepository) GetAllSubjects(ctx context.Context) ([]entity.Subject, error) {
	allSubjectsCacheKey := fmt.Sprintf("%sall", cache.SUBJECT_CACHE_KEY)

	cachedSubjects, err := r.redisClient.Get(ctx, allSubjectsCacheKey).Result()
	if err == nil && cachedSubjects != "" {
		var subjects []entity.Subject
		if json.Unmarshal([]byte(cachedSubjects), &subjects) == nil {
			return subjects, nil
		}
	}

	subjects, err := r.client.Subject.FindMany().With(db.Subject.User.Fetch()).Exec(ctx)
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

	subjectsJSON, _ := json.Marshal(result)
	r.redisClient.Set(ctx, allSubjectsCacheKey, string(subjectsJSON), time.Duration(cache.SUBJECT_CACHE_KEY_TTL)*time.Second)

	return result, nil
}

// NOTE - get subject by id repository
func (r *subjectRepository) GetSubjectByID(ctx context.Context, id int) (*entity.Subject, error) {
	subjectCacheKey := fmt.Sprintf("%s%d", cache.SUBJECT_CACHE_KEY, id)

	cachedSubject, err := r.redisClient.Get(ctx, subjectCacheKey).Result()
	if err == nil && cachedSubject != "" {
		var subject entity.Subject
		if json.Unmarshal([]byte(cachedSubject), &subject) == nil {
			return &subject, nil
		}
	}

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

	subjectData, _ := json.Marshal(subject)
	r.redisClient.Set(ctx, subjectCacheKey, string(subjectData), time.Duration(cache.SUBJECT_CACHE_KEY_TTL)*time.Second)

	return &entity.Subject{
		ID:        subject.ID,
		Name:      subject.Name,
		User:      users,
		Status:    subject.Status,
		CreatedAt: utils.FormatToVientianeTime(subject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(subject.UpdatedAt),
	}, nil
}

// NOTE - create subject repository
func (r *subjectRepository) CreateSubject(ctx context.Context, subject entity.Subject) (*entity.Subject, error) {
	newSubject, err := r.client.Subject.CreateOne(
		db.Subject.Name.Set(subject.Name),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	r.redisClient.Del(ctx, fmt.Sprintf("%sall", cache.SUBJECT_CACHE_KEY))

	return &entity.Subject{
		ID:        newSubject.ID,
		Name:      newSubject.Name,
		Status:    newSubject.Status,
		CreatedAt: utils.FormatToVientianeTime(newSubject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(newSubject.UpdatedAt),
	}, nil
}

// NOTE - update subject repository
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

	subjectCacheKey := fmt.Sprintf("%s%d", cache.SUBJECT_CACHE_KEY, id)
	r.redisClient.Del(ctx, subjectCacheKey)
	r.redisClient.Del(ctx, fmt.Sprintf("%sall", cache.SUBJECT_CACHE_KEY))

	return &entity.Subject{
		ID:        updateSubject.ID,
		Name:      updateSubject.Name,
		Status:    updateSubject.Status,
		CreatedAt: utils.FormatToVientianeTime(updateSubject.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(updateSubject.UpdatedAt),
	}, nil
}

// NOTE - delete subject repository
func (r *subjectRepository) DeleteSubject(ctx context.Context, id int) error {
	_, err := r.client.Subject.FindUnique(
		db.Subject.ID.Equals(id),
	).Delete().Exec(ctx)

	subjectCacheKey := fmt.Sprintf("%s%d", cache.SUBJECT_CACHE_KEY, id)
	r.redisClient.Del(ctx, subjectCacheKey)
	r.redisClient.Del(ctx, fmt.Sprintf("%sall", cache.SUBJECT_CACHE_KEY))

	return err
}

// NOTE - clear subject cache repository
func (r *subjectRepository) ClearSubjectCache(ctx context.Context) error {
	pattern := fmt.Sprintf("%s*", cache.SUBJECT_CACHE_KEY)
	iter := r.redisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		if err := r.redisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
