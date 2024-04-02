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
	var params types.AddTaskRequest
	id := c.Params("id")
	if len(id) == 0 {
		return ErrBadRequest()
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	project, err := h.getProjectByID(c.Context(), id)
	if err != nil {
		return err
	}
	task, err := h.getTaskByID(c.Context(), params.TaskID)
	if err != nil {
		return err
	}
	//TODO: Check if the task already exists on this project
	if err := h.addTaskToProject(c.Context(), project.ID, task.ID); err != nil {
		return err
	}
	if err := h.updateTaskProjectID(c.Context(), project.ID, task.ID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": project.ID})
}
func (h *ProjectHandler) getProjectByID(ctx context.Context, id string) (*types.Project, error) {
	project, err := h.store.Project.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("project")
		}
		return nil, err
	}
	return project, nil
}
func (h *ProjectHandler) getTaskByID(ctx context.Context, id string) (*types.Task, error) {
	task, err := h.store.Task.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return nil, ErrResourceNotFound("task")
		}
		return nil, err
	}
	return task, nil
}
func (h *ProjectHandler) addTaskToProject(ctx context.Context, projectID, taskID string) error {
	filter := db.NewMap("_id", projectID)
	update := db.PushToKey("tasks", taskID)
	if err := h.store.Project.Update(ctx, filter, update); err != nil {
		return err
	}
	return nil
}
func (h *ProjectHandler) updateTaskProjectID(ctx context.Context, projectID, taskID string) error {
	filter := db.NewMap("_id", taskID)
	update := db.PushToKey("projects", projectID)
	return h.store.Task.Update(ctx, filter, update)
}
