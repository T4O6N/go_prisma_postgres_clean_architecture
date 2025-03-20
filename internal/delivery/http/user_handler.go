package delivery

import (
	"context"
	"net/http"
	"sample-project/internal/entity"
	"sample-project/internal/usecase"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// NOTE - user handler struct
type UserHandler struct {
	useCase usecase.UserUseCase
}

// NOTE - new user handler
func NewUserHandler(router *gin.Engine, useCase usecase.UserUseCase) {
	handler := &UserHandler{useCase: useCase}

	users := router.Group("/api/v1/users")

	users.GET("", handler.GetUsers)
	users.GET("/:id", handler.GetUserByID)
	users.POST("", handler.CreateUser)
	users.PUT("/update/:id", handler.UpdateUser)
	users.DELETE("/delete/:id", handler.DeleteUser)
	users.DELETE("/clear-cache", handler.ClearUserCache)
}

// NOTE - get all users handler
// @Summary Get all users
// @Description Get list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Results per page (default: 10)" minimum(1) maximum(100)
// @Param name query string false "Filter by user name (partial match)"
// @Param startDate query string false "Filter by start date (format: YYYY-MM-DD)"
// @Param endDate query string false "Filter by end date (format: YYYY-MM-DD)"
// @Success 200 {object} entity.UserListResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}

	name := c.DefaultQuery("name", "")
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	users, totalCount, err := h.useCase.GetUsers(c, page, limit, name, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      totalCount,
			"totalPages": (totalCount + limit - 1) / limit,
		},
		"data": users,
	})
}

// NOTE - get user by id handler
// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.UserResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.useCase.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// NOTE - create user handler
// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.CreateUserRequest true "User data"
// @Success 201 {object} entity.UserResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req entity.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := entity.User{
		Name:      req.Name,
		Email:     req.Email,
		SubjectID: req.SubjectID,
		Status:    true,
	}

	createdUser, err := h.useCase.CreateUser(c, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

// NOTE - update user handler
// @Summary Update a user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body entity.UpdateUserRequest true "User data"
// @Success 201 {object} entity.UserResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users/update/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	_, err = h.useCase.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req entity.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUpdate := entity.User{
		Name:      req.Name,
		Email:     req.Email,
		Status:    req.Status,
		SubjectID: req.SubjectID,
	}

	updatedUser, err := h.useCase.UpdateUser(c, id, userUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// NOTE - delete user handler
// @Summary Delete a user
// @Description Remove a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} entity.ErrorResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users/delete/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.useCase.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// NOTE - clear cache of users handler
// @Summary Clear cache of users
// @Description Clear the cache of users
// @Tags users
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 404 {object} entity.ErrorResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/users/clear-cache [delete]
func (h *UserHandler) ClearUserCache(c *gin.Context) {
	ctx := context.Background()
	err := h.useCase.ClearUserCache(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear users caches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User caches cleared successfully"})
}
