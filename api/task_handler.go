package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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
	task, err := h.taskStore.GetTaskByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound()
		}
		return err
	}
	return c.JSON(task)
}

func (h *TaskHandler) HandleGetTasks(c *fiber.Ctx) error {
	tasks, err := h.taskStore.GetTasks(c.Context())
	if err != nil {
		return ErrResourceNotFound()
	}
	return c.JSON(tasks)
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
		return NewError(http.StatusBadRequest, "Task already completed")
	}
	filter := db.Map{"_id": id}
	params := types.UpdateTaskParams{
		Completed: true,
	}
	if err = h.taskStore.UpdateTask(c.Context(), filter, params); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound()
		}
		return err
	}
	return c.JSON(map[string]string{"updated": id})

}
func (h *TaskHandler) HandleDeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.taskStore.DeleteTask(c.Context(), id); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound()
		}
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *TaskHandler) isTaskCompleted(ctx context.Context, id string) (bool, error) {
	task, err := h.taskStore.GetTaskByID(ctx, id)
	if err != nil {
		return false, err
	}
	return task.Completed, nil
}
