package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskStore db.TaskStore
}

func NewTaskHandler(taskStore db.TaskStore) *TaskHandler {
	return &TaskHandler{
		taskStore: taskStore,
	}
}
func (h *TaskHandler) HandleGetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := h.taskStore.GetTaskByID(c.Context(), types.ID(id))
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("task")
		}
		return err
	}
	return c.JSON(task)
}

type TaskQueryParams struct {
	db.Pagination
	Completed *bool
}

// TODO: review
func createCompletionFilter(completed *bool) db.Map {
	if completed != nil {
		return db.Map{"completed": completed}
	} else {
		return nil
	}
}
func (h *TaskHandler) HandleGetTasks(c *fiber.Ctx) error {
	var params TaskQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	tasks, err := h.taskStore.GetTasks(c.Context(), createCompletionFilter(params.Completed), &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("task")
	}
	resp := NewResourceResponse(tasks, len(tasks), params.Page)
	return c.JSON(resp)
}

func (h *TaskHandler) HandlePostTask(c *fiber.Ctx) error {
	var params types.CreateTaskParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	task := types.NewTaskFromParams(params)
	insertedTask, err := h.taskStore.InsertTask(c.Context(), task)
	if err != nil {
		return err
	}
	return c.JSON(insertedTask)
}

func (h *TaskHandler) HandleCompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	completed, err := h.isTaskCompleted(c.Context(), id)
	if err != nil {
		return err
	}
	if completed {
		return ErrBadRequestCustomMessage("task already completed")
	}
	params := types.UpdateTaskParams{
		Completed: true,
	}
	update := db.SetUpdateMap(params.ToMap())
	if err = h.taskStore.Update(c.Context(), types.ID(id), update); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("task")
		}
		return err
	}
	return c.JSON(fiber.Map{"updated": id})

}
func (h *TaskHandler) HandleDeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.taskStore.Delete(c.Context(), types.ID(id)); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("task")
		}
		return err
	}
	return c.JSON(fiber.Map{"deleted": id})
}

func (h *TaskHandler) isTaskCompleted(ctx context.Context, id string) (bool, error) {
	task, err := h.taskStore.GetTaskByID(ctx, types.ID(id))
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return false, ErrResourceNotFound("task")
		}
		return false, err
	}
	return task.Completed, nil
}
