package v1

import (
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type studentRoutes struct {
	s usecase.Student
	l logger.Interface
	v *validator.Validate
}

func NewStudentRoutes(router fiber.Router, s usecase.Student, l logger.Interface) {
	r := &studentRoutes{s, l, validator.New(validator.WithRequiredStructEnabled())}

	// Register routes
	router.Post("/students", r.createStudent)
	router.Get("/students", r.getStudents)
	router.Get("/students/:id", r.getStudentByID)
	router.Put("/students/:id", r.updateStudent)
	router.Delete("/students/:id", r.deleteStudent)
}

type createStudentRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	GroupID int    `json:"group_id" validate:"required"`
}

// @Summary     Create a student
// @Description Add a new student to the system
// @ID          create-student
// @Tags  	    students
// @Accept      json
// @Produce     json
// @Param       request body createStudentRequest true "Student data"
// @Success     201 {object} entity.Student
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /students [post]
func (r *studentRoutes) createStudent(ctx *fiber.Ctx) error {
	var request createStudentRequest

	if err := ctx.BodyParser(&request); err != nil {
		r.l.Error(err, "http - v1 - createStudent")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Error(err, "http - v1 - createStudent - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	student := entity.Student{
		Name:    request.Name,
		Email:   request.Email,
		GroupID: request.GroupID,
	}

	createdStudent, err := r.s.CreateStudent(ctx.UserContext(), student)
	if err != nil {
		r.l.Error(err, "http - v1 - createStudent - r.s.CreateStudent")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create student")
	}

	return ctx.Status(http.StatusCreated).JSON(createdStudent)
}

// @Summary     Get all students
// @Description Retrieve a list of all students
// @ID          get-students
// @Tags  	    students
// @Accept      json
// @Produce     json
// @Param       query query string false "Search query"
// @Success     200 {array} entity.Student
// @Failure     500 {object} response
// @Router      /students [get]
func (r *studentRoutes) getStudents(ctx *fiber.Ctx) error {
	// Check if there's a search query
	query := ctx.Query("query")
	if query != "" {
		return r.searchStudents(ctx, query)
	}

	students, err := r.s.GetStudents(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getStudents - r.s.GetStudents")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get students")
	}

	return ctx.Status(http.StatusOK).JSON(students)
}

// Search students based on query
func (r *studentRoutes) searchStudents(ctx *fiber.Ctx, query string) error {
	students, err := r.s.SearchStudents(ctx.UserContext(), query)
	if err != nil {
		r.l.Error(err, "http - v1 - searchStudents - r.s.SearchStudents")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to search students")
	}

	return ctx.Status(http.StatusOK).JSON(students)
}

// @Summary     Get student by ID
// @Description Retrieve a specific student by ID
// @ID          get-student-by-id
// @Tags  	    students
// @Accept      json
// @Produce     json
// @Param       id path int true "Student ID"
// @Success     200 {object} entity.Student
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /students/{id} [get]
func (r *studentRoutes) getStudentByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - getStudentByID")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	student, err := r.s.GetStudentByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getStudentByID - r.s.GetStudentByID")
		return errorResponse(ctx, http.StatusNotFound, "student not found")
	}

	return ctx.Status(http.StatusOK).JSON(student)
}

type updateStudentRequest struct {
	Name    string `json:"name" validate:"required"`
	GroupID int    `json:"group_id" validate:"required"`
}

// @Summary     Update student
// @Description Update a student's information
// @ID          update-student
// @Tags  	    students
// @Accept      json
// @Produce     json
// @Param       id path int true "Student ID"
// @Param       request body updateStudentRequest true "Updated student data"
// @Success     200 {object} entity.Student
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /students/{id} [put]
func (r *studentRoutes) updateStudent(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - updateStudent")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	var request updateStudentRequest
	if err := ctx.BodyParser(&request); err != nil {
		r.l.Error(err, "http - v1 - updateStudent")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Error(err, "http - v1 - updateStudent - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	// First check if student exists
	_, err = r.s.GetStudentByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - updateStudent - r.s.GetStudentByID")
		return errorResponse(ctx, http.StatusNotFound, "student not found")
	}

	student := entity.Student{
		ID:      id,
		Name:    request.Name,
		GroupID: request.GroupID,
	}

	err = r.s.UpdateStudent(ctx.UserContext(), student)
	if err != nil {
		r.l.Error(err, "http - v1 - updateStudent - r.s.UpdateStudent")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update student")
	}

	// Get the updated student to return in response
	updatedStudent, err := r.s.GetStudentByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - updateStudent - r.s.GetStudentByID")
		return errorResponse(ctx, http.StatusInternalServerError, "student updated but failed to retrieve updated data")
	}

	return ctx.Status(http.StatusOK).JSON(updatedStudent)
}

// @Summary     Delete student
// @Description Remove a student from the system
// @ID          delete-student
// @Tags  	    students
// @Accept      json
// @Produce     json
// @Param       id path int true "Student ID"
// @Success     204 "No Content"
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /students/{id} [delete]
func (r *studentRoutes) deleteStudent(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - deleteStudent")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	// First check if student exists
	_, err = r.s.GetStudentByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteStudent - r.s.GetStudentByID")
		return errorResponse(ctx, http.StatusNotFound, "student not found")
	}

	err = r.s.DeleteStudent(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteStudent - r.s.DeleteStudent")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete student")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
