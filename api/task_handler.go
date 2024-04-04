package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/data"
	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/service"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}
func (h *TaskHandler) HandleGetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := h.taskService.GetTaskByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("task")
		}
		return err
	}
	return c.JSON(task)
}

func (h *TaskHandler) HandleGetTasks(c *fiber.Ctx) error {
	var params db.TaskQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	tasks, err := h.taskService.GetTasks(c.Context(), &params)
	if err != nil {
		return ErrResourceNotFound("task")
	}
	resp := NewResourceResponse(tasks, len(tasks), params.Page)
	return c.JSON(resp)
}

func (h *TaskHandler) HandlePostTask(c *fiber.Ctx) error {
	var params data.CreateTaskParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	task, err := h.taskService.CreateTask(c.Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(task)
}

func (h *TaskHandler) HandleCompleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.taskService.CompleteTask(c.Context(), id); err != nil {
		switch {
		case errors.Is(err, db.ErrorNotFound):
			return ErrResourceNotFound("task")
		case errors.Is(err, service.ErrTaskAlreadyCompleted):
			return ErrBadRequestCustomMessage(err.Error())
		default:
			return err
		}
	}
	return c.JSON(fiber.Map{"updated": id})

}
func (h *TaskHandler) HandleDeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.taskService.DeleteTask(c.Context(), id); err != nil {
		if errors.Is(err, db.ErrorNotFound) {
			return ErrResourceNotFound("task")
		}
		return err
	}
	return c.JSON(fiber.Map{"deleted": id})
}
