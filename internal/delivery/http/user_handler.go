package delivery

import (
	"context"
	"net/http"
	"sample-project/internal/dto"
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
// @Success 200 {array} entity.User
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.useCase.GetUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// NOTE - get user by id handler
// @Summary Get user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.User
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
// @Param user body dto.CreateUserDto true "User data"
// @Success 201 {object} entity.User
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserDto dto.CreateUserDto
	if err := c.ShouldBindJSON(&createUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entity.User{
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		SubjectID: createUserDto.SubjectID,
	}

	newUser, err := h.useCase.CreateUser(c, user)
	if err != nil {
		// Check if the error is related to "not found"
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// For other errors, return internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

// NOTE - update user handler
// @Summary Update a user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserDto true "Updated user data"
// @Success 200 {object} entity.User
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

	var updateUserDto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&updateUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.useCase.UpdateUser(c, id, updateUserDto)
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
