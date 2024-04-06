package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/service"
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
	var params data.CreateProjectParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrInternalServer()
	}
	insertedProject, err := h.projectService.CreateProject(c.Context(), params, user.ID)
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
	project, err := h.projectService.GetProjectByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			return ErrResourceNotFound("project")
		}
		return err
	}
	return c.JSON(project)
}
