package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	store *db.Store
}

func NewProjectHandler(store *db.Store) *ProjectHandler {
	return &ProjectHandler{
		store: store,
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
	project := data.NewProjectFromParams(params)
	project.UserID = user.ID
	insertedProject, err := h.store.Project.InsertProject(c.Context(), project)
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
	project, err := h.store.Project.GetProjectByID(c.Context(), data.ID(id))
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
	tasks, err := h.store.Task.GetTasksByProject(c.Context(), data.ID(id), pagination)
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
	project, err := h.getProjectByID(c.Context(), data.ID(id))
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if user.ID != project.UserID {
		return ErrUnAuthorized()
	}
	task, err := h.getTaskByID(c.Context(), data.ID(params.TaskID))
	if err != nil {
		return err
	}
	if project.ContainsTask(task.ID) {
		return ErrBadRequestCustomMessage(fmt.Sprintf("task %s is already associated with this project", task.ID))
	}
	if err := h.store.Project.UpdateProjectTasks(c.Context(), data.ID(id), data.ID(params.TaskID)); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": data.ID(id)})
}
func (h *ProjectHandler) getProjectByID(ctx context.Context, id data.ID) (*data.Project, error) {
	project, err := h.store.Project.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("project")
		}
		return nil, err
	}
	return project, nil
}
func (h *ProjectHandler) getTaskByID(ctx context.Context, id data.ID) (*data.Task, error) {
	task, err := h.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("task")
		}
		return nil, err
	}
	return task, nil
}
