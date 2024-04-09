package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) HandlePostProject(c *fiber.Ctx) error {
	var params types.NewProjectParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	auth, err := getAuth(c)
	if err != nil {
		return ErrInternalServer()
	}
	insertedProject, err := h.projectService.CreateProject(c.Context(), params, auth.UserID)
	if err != nil {
		return err
	}
	return c.JSON(insertedProject)

}
func (h *ProjectHandler) HandleGetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	auth, err := getAuth(c)
	if err != nil {
		return err
	}
	project, err := h.projectService.GetProjectByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			return ErrResourceNotFound(err.Error())
		}
		return err
	}

	if project.UserID != auth.UserID {
		return ErrUnAuthorized()
	}
	return c.JSON(project)
}
