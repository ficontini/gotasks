package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
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
	var params types.CreateProjectParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrInternalServer()
	}
	project := types.NewProjectFromParams(params)
	project.UserID = user.ID
	insertedProject, err := h.store.Project.InsertProject(c.Context(), project)
	if err != nil {
		return err
	}
	return c.JSON(insertedProject)

}
func (h *ProjectHandler) HandleGetTasks(c *fiber.Ctx) error {

	return nil
}
func (h *ProjectHandler) HandleAddTaskToProject(c *fiber.Ctx) error {
	var params types.AddTaskParams
	id := c.Params("id")
	if len(id) == 0 {
		return ErrBadRequest()
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	project, err := h.getProjectByID(c.Context(), types.ID(id))
	if err != nil {
		return err
	}
	task, err := h.getTaskByID(c.Context(), types.ID(params.TaskID))
	if err != nil {
		return err
	}
	//TODO: Check if the task already exists on this project
	if err := h.store.Project.Update(c.Context(), project.ID, db.PushToKey(params.ToMap())); err != nil {
		return err
	}

	if err := h.store.Task.Update(c.Context(), task.ID, db.PushToKey(db.Map{"projects": project.ID})); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": project.ID})
}
func (h *ProjectHandler) getProjectByID(ctx context.Context, id types.ID) (*types.Project, error) {
	project, err := h.store.Project.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("project")
		}
		return nil, err
	}
	return project, nil
}
func (h *ProjectHandler) getTaskByID(ctx context.Context, id types.ID) (*types.Task, error) {
	task, err := h.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("task")
		}
		return nil, err
	}
	return task, nil
}
