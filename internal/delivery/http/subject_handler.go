package delivery

import (
	"net/http"
	"sample-project/internal/dto"
	"sample-project/internal/entity"
	"sample-project/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	useCase usecase.SubjectUsecase
}

func NewSubjectHandler(router *gin.Engine, useCase usecase.SubjectUsecase) {
	handler := &SubjectHandler{useCase: useCase}

	subjects := router.Group("/api/v1/subjects")

	subjects.GET("", handler.GetSubject)
	subjects.GET("/:id", handler.GetSubjectByID)
	subjects.POST("", handler.CreateSubject)
	subjects.PUT("/update/:id", handler.UpdateSubject)
	subjects.DELETE("/delete/:id", handler.DeleteSubject)
}

// @Summary Get all subjects
// @Description Get a list of all subjects
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {array} entity.Subject
// @Router /api/v1/subjects [get]
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	subjects, err := h.useCase.GetSubject(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subjects)
}

// @Summary Get subject by ID
// @Description Get a single subject by ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} entity.Subject
// @Router /api/v1/subjects/{id} [get]
func (h *SubjectHandler) GetSubjectByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}

	subject, err := h.useCase.GetSubjectByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// @Summary Create a subject
// @Description Create a new subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param subject body dto.CreateSubjectDto true "Subject data"
// @Success 201 {object} entity.Subject
// @Router /api/v1/subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var createSubjectDto dto.CreateSubjectDto
	if err := c.ShouldBindJSON(&createSubjectDto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subject := entity.Subject{
		Name: createSubjectDto.Name,
	}

	newSubject, err := h.useCase.CreateSubject(c, subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newSubject)
}

// @Summary Update a subject
// @Description Update subject details
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param subject body dto.UpdateSubjectDto true "Updated subject data"
// @Success 200 {object} entity.Subject
// @Router /api/v1/subjects/update/{id} [put]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}

	var updateSubjectDto dto.UpdateSubjectDto
	if err := c.ShouldBindJSON(&updateSubjectDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSubject, err := h.useCase.UpdateSubject(c, id, updateSubjectDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSubject)
}

// @Summary Delete a subject
// @Description Remove a subject by ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Success 204 "No Content"
// @Router /api/v1/subjects/delete/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}

	err = h.useCase.DeleteSubject(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
