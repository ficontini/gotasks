package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
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

func (h *ProjectHandler) HandleGetTasks(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	project, err := h.projectService.GetProjectByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("project")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if project.UserID != user.ID {
		return ErrUnAuthorized()
	}
	//TODO: implement pagination params
	pagination := &db.Pagination{}
	tasks, err := h.projectService.GetTasksByProject(c.Context(), id, pagination)
	if err != nil {
		return err
	}
	resp := NewResourceResponse(tasks, len(tasks), pagination.Page)
	return c.JSON(resp)
}
func (h *ProjectHandler) HandleAddTaskToProject(c *fiber.Ctx) error {
	var params data.AddTaskParams
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	project, err := h.projectService.GetProjectByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("project")
		}
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if user.ID != project.UserID {
		return ErrUnAuthorized()
	}
	err = h.projectService.AddTask(c.Context(), project, params.TaskID)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrorNotFound):
			return ErrResourceNotFound("task")
		case errors.Is(err, service.ErrTaskAlreadyAssociated):
			return ErrBadRequestCustomMessage(err.Error())
		default:
			return err
		}
	}
	return c.JSON(fiber.Map{"updated": data.ID(id)})
}
