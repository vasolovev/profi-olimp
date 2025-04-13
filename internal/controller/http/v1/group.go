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

type groupRoutes struct {
	g usecase.Group
	l logger.Interface
	v *validator.Validate
}

func NewGroupRoutes(router fiber.Router, g usecase.Group, l logger.Interface) {
	r := &groupRoutes{g, l, validator.New(validator.WithRequiredStructEnabled())}

	// Register routes
	router.Post("/groups", r.createGroup)
	router.Get("/groups", r.getGroups)
	router.Get("/groups/:id", r.getGroupByID)
	router.Put("/groups/:id", r.updateGroup)
	router.Delete("/groups/:id", r.deleteGroup)
}

type createGroupRequest struct {
	Name     string `json:"name" validate:"required"`
	ParentID *int   `json:"parent_id"`
}

// @Summary     Create a group
// @Description Add a new academic group
// @ID          create-group
// @Tags  	    groups
// @Accept      json
// @Produce     json
// @Param       request body createGroupRequest true "Group data"
// @Success     201 {object} entity.Group
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /groups [post]
func (r *groupRoutes) createGroup(ctx *fiber.Ctx) error {
	var request createGroupRequest

	if err := ctx.BodyParser(&request); err != nil {
		r.l.Error(err, "http - v1 - createGroup")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Error(err, "http - v1 - createGroup - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	group := entity.Group{
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	createdGroup, err := r.g.CreateGroup(ctx.UserContext(), group)
	if err != nil {
		r.l.Error(err, "http - v1 - createGroup - r.g.CreateGroup")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create group")
	}

	return ctx.Status(http.StatusCreated).JSON(createdGroup)
}

// @Summary     Get all groups
// @Description Retrieve a list of all academic groups with their subgroups
// @ID          get-groups
// @Tags  	    groups
// @Accept      json
// @Produce     json
// @Param       query query string false "Search query"
// @Success     200 {array} entity.Group
// @Failure     500 {object} response
// @Router      /groups [get]
func (r *groupRoutes) getGroups(ctx *fiber.Ctx) error {
	// Check if there's a search query
	query := ctx.Query("query")
	if query != "" {
		return r.searchGroups(ctx, query)
	}

	groups, err := r.g.GetGroups(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getGroups - r.g.GetGroups")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get groups")
	}

	return ctx.Status(http.StatusOK).JSON(groups)
}

// Search groups based on query
func (r *groupRoutes) searchGroups(ctx *fiber.Ctx, query string) error {
	groups, err := r.g.SearchGroups(ctx.UserContext(), query)
	if err != nil {
		r.l.Error(err, "http - v1 - searchGroups - r.g.SearchGroups")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to search groups")
	}

	return ctx.Status(http.StatusOK).JSON(groups)
}

// @Summary     Get group by ID
// @Description Retrieve a specific academic group by ID
// @ID          get-group-by-id
// @Tags  	    groups
// @Accept      json
// @Produce     json
// @Param       id path int true "Group ID"
// @Success     200 {object} entity.Group
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /groups/{id} [get]
func (r *groupRoutes) getGroupByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - getGroupByID")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	group, err := r.g.GetGroupByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getGroupByID - r.g.GetGroupByID")
		return errorResponse(ctx, http.StatusNotFound, "group not found")
	}

	return ctx.Status(http.StatusOK).JSON(group)
}

type updateGroupRequest struct {
	Name     string `json:"name" validate:"required"`
	ParentID *int   `json:"parent_id"`
}

// @Summary     Update group
// @Description Update an academic group's information
// @ID          update-group
// @Tags  	    groups
// @Accept      json
// @Produce     json
// @Param       id path int true "Group ID"
// @Param       request body updateGroupRequest true "Updated group data"
// @Success     200 {object} entity.Group
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /groups/{id} [put]
func (r *groupRoutes) updateGroup(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - updateGroup")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	var request updateGroupRequest
	if err := ctx.BodyParser(&request); err != nil {
		r.l.Error(err, "http - v1 - updateGroup")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Error(err, "http - v1 - updateGroup - validation")
		return errorResponse(ctx, http.StatusBadRequest, "validation failed")
	}

	// First check if group exists
	_, err = r.g.GetGroupByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - updateGroup - r.g.GetGroupByID")
		return errorResponse(ctx, http.StatusNotFound, "group not found")
	}

	group := entity.Group{
		ID:       id,
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	err = r.g.UpdateGroup(ctx.UserContext(), group)
	if err != nil {
		r.l.Error(err, "http - v1 - updateGroup - r.g.UpdateGroup")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update group")
	}

	// Get the updated group to return in response
	updatedGroup, err := r.g.GetGroupByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - updateGroup - r.g.GetGroupByID")
		return errorResponse(ctx, http.StatusInternalServerError, "group updated but failed to retrieve updated data")
	}

	return ctx.Status(http.StatusOK).JSON(updatedGroup)
}

// @Summary     Delete group
// @Description Remove an academic group from the system
// @ID          delete-group
// @Tags  	    groups
// @Accept      json
// @Produce     json
// @Param       id path int true "Group ID"
// @Success     204 "No Content"
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     409 {object} response
// @Failure     500 {object} response
// @Router      /groups/{id} [delete]
func (r *groupRoutes) deleteGroup(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - deleteGroup")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
	}

	// First check if group exists
	_, err = r.g.GetGroupByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteGroup - r.g.GetGroupByID")
		return errorResponse(ctx, http.StatusNotFound, "group not found")
	}

	err = r.g.DeleteGroup(ctx.UserContext(), id)
	if err != nil {
		if err.Error() == "cannot delete group with subgroups" {
			return errorResponse(ctx, http.StatusConflict, "cannot delete group with subgroups")
		}
		r.l.Error(err, "http - v1 - deleteGroup - r.g.DeleteGroup")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete group")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
