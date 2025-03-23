package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sample-project/internal/config/cache"
	"sample-project/internal/entity"
	"sample-project/internal/utils"
	"sample-project/prisma/db"
	"time"

	"github.com/redis/go-redis/v9"
)

// NOTE - user repository interface
type UserRepository interface {
	GetAllUsers(ctx context.Context, page, limit int, name string, startDate, endDate string) ([]entity.User, int, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	ClearUserCache(ctx context.Context) error
}

// NOTE - user repository struct
type userRepository struct {
	client      *db.PrismaClient
	redisClient *redis.Client
}

// NOTE - new user repository
func NewUserRepository(client *db.PrismaClient, redisClient *redis.Client) UserRepository {
	return &userRepository{client: client, redisClient: redisClient}
}

// NOTE - get all users repository
func (r *userRepository) GetAllUsers(ctx context.Context, page, limit int, name string, startDate, endDate string) ([]entity.User, int, error) {
	offset := (page - 1) * limit
	allUsersCacheKey := fmt.Sprintf("%sall_page%d_limit%d_name%s_start%s_end%s", cache.USER_CACHE_KEY, page, limit, name, startDate, endDate)

	// Check Redis Cache First
	cachedUsers, err := r.redisClient.Get(ctx, allUsersCacheKey).Result()
	if err == nil && cachedUsers != "" {
		var users []entity.User
		if json.Unmarshal([]byte(cachedUsers), &users) == nil {
			return users, len(users), nil
		}
	}

	whereClause := []db.UserWhereParam{}
	if name != "" {
		whereClause = append(whereClause, db.User.Name.Contains(name))
	}

	if startDate != "" {
		startTime, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			// Ensure start time is the beginning of the day (00:00:00)
			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
			whereClause = append(whereClause, db.User.CreatedAt.Gte(startTime))
		}
	}

	if endDate != "" {
		endTime, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// Ensure end time is the last second of the day (23:59:59)
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.UTC)
			whereClause = append(whereClause, db.User.CreatedAt.Lte(endTime))
		}
	}

	// Fetch users from DB
	users, err := r.client.User.FindMany(whereClause...).
		Skip(offset).
		Take(limit).
		OrderBy(db.User.CreatedAt.Order(db.SortOrderDesc)).
		Exec(ctx)
	if err != nil {
		return nil, 0, err
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
			Password:  u.Password,
			SubjectID: subjectID,
			Status:    u.Status,
			CreatedAt: utils.FormatToVientianeTime(u.CreatedAt),
			UpdatedAt: utils.FormatToVientianeTime(u.UpdatedAt),
		})
	}

	// Store in Redis Cache
	usersJSON, _ := json.Marshal(result)
	r.redisClient.Set(ctx, allUsersCacheKey, string(usersJSON), time.Duration(cache.USER_CACHE_KEY_TTL)*time.Second)

	return result, len(result), nil
}

// NOTE - get user by id repository
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	userCacheKey := fmt.Sprintf("%s%d", cache.USER_CACHE_KEY, id)

	// Check if user exists in cache
	cachedUser, err := r.redisClient.Get(ctx, userCacheKey).Result()
	if err == nil && cachedUser != "" {
		var user entity.User
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil
		}
	}

	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Store in Redis cache with TTL
	userData, _ := json.Marshal(user)
	r.redisClient.Set(ctx, userCacheKey, string(userData), time.Duration(cache.USER_CACHE_KEY_TTL))

	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Status:    user.Status,
		CreatedAt: utils.FormatToVientianeTime(user.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(user.UpdatedAt),
	}, nil
}

// NOTE - get user by email
func (r *userRepository) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	user, err := r.client.User.FindFirst(
		db.User.Name.Equals(name),
	).Exec(ctx)
	if err != nil {
		slog.Error("Failed to fetch user by name", "name", name, "error", err)
		return nil, err
	}

	slog.Info("Fetched user by name", "name", name, "user", user)
	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Status:    user.Status,
		CreatedAt: utils.FormatToVientianeTime(user.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(user.UpdatedAt),
	}, nil
}

// NOTE - create user repository
func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

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

	// Get current time in Vientiane timezone
	currentTime := utils.FormatToVientianeTime(time.Now())

	// Extract day, month, and year from creation time
	day := currentTime.Day()
	month := int(currentTime.Month())
	year := currentTime.Year()

	newUser, err := r.client.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(hashedPassword),
		db.User.Day.Set(day),
		db.User.Month.Set(month),
		db.User.Year.Set(year),
		db.User.SubjectID.Set(user.SubjectID),
		db.User.Status.Set(user.Status),
		db.User.CreatedAt.Set(currentTime),
		db.User.UpdatedAt.Set(currentTime),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Clear cache after create
	cache.DelWithPattern(ctx, fmt.Sprintf("%sall*", cache.USER_CACHE_KEY))

	return &entity.User{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Password:  hashedPassword,
		SubjectID: user.SubjectID,
		Status:    user.Status,
		Day:       newUser.Day,
		Month:     newUser.Month,
		Year:      newUser.Year,
		CreatedAt: utils.FormatToVientianeTime(newUser.CreatedAt),
		UpdatedAt: utils.FormatToVientianeTime(newUser.UpdatedAt),
	}, nil
}

// NOTE - update user repository
func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (*entity.User, error) {
	var updates []db.UserSetParam

	updates = append(updates, db.User.Name.Set(user.Name))
	updates = append(updates, db.User.Email.Set(user.Email))
	updates = append(updates, db.User.Status.Set(user.Status))
	updates = append(updates, db.User.UpdatedAt.Set(utils.FormatToVientianeTime(time.Now())))

	if user.SubjectID != 0 {
		subject, err := r.client.Subject.FindUnique(
			db.Subject.ID.Equals(user.SubjectID),
		).Exec(ctx)

		if err != nil || subject == nil {
			return nil, fmt.Errorf("subject with ID: %d not found", user.SubjectID)
		}

		updates = append(updates, db.User.SubjectID.Set(user.SubjectID))
	}

	updateUser, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		updates...,
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// clear cache after updating
	userCacheKey := fmt.Sprintf("%s%d", cache.USER_CACHE_KEY, id)
	cache.Del(ctx, userCacheKey)
	cache.DelWithPattern(ctx, fmt.Sprintf("%sall*", cache.USER_CACHE_KEY))

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

// NOTE - delete user repository
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(ctx)

	userCacheKey := fmt.Sprintf("%s%d", cache.USER_CACHE_KEY, id)
	cache.Del(ctx, userCacheKey)
	cache.DelWithPattern(ctx, fmt.Sprintf("%sall*", cache.USER_CACHE_KEY))

	return err
}

// NOTE - clear user cache repository
func (r *userRepository) ClearUserCache(ctx context.Context) error {
	pattern := fmt.Sprintf("%s*", cache.USER_CACHE_KEY)
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
