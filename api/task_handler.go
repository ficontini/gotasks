package api

import (
	"errors"
	"net/http"

	"github.com/ficontini/gotasks/service"
	"github.com/ficontini/gotasks/types"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskService service.TaskServicer
}

func NewTaskHandler(taskService service.TaskServicer) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}
func (h *TaskHandler) HandleGetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	task, err := h.taskService.GetTaskByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return ErrResourceNotFound(err.Error())
		}
		return err
	}
	return c.JSON(task)
}
func (h *TaskHandler) HandleGetUserTasks(c *fiber.Ctx) error {
	auth, err := getAuth(c)
	if err != nil {
		return err
	}
	var params service.TaskQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	tasks, err := h.taskService.GetTasksByUserID(c.Context(), auth.UserID, params)
	if err != nil {
		return err
	}
	resp := NewResourceResponse(tasks, len(tasks), params.Page)
	return c.JSON(resp)
}
func (h *TaskHandler) HandleGetTasks(c *fiber.Ctx) error {
	var params service.TaskQueryParams
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
	var params types.NewTaskParams
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
	if len(id) == 0 {
		return ErrInvalidID()
	}
	user, err := getUserAuth(c)
	if err != nil {
		return err
	}
	params := types.UpdateTaskRequest{
		TaskID: id,
		UserID: user.ID,
	}
	if err := h.taskService.CompleteTask(c.Context(), params); err != nil {
		switch {
		case errors.Is(err, service.ErrUnAuthorized):
			return ErrUnAuthorized()
		case errors.Is(err, service.ErrTaskNotFound):
			return ErrResourceNotFound(err.Error())
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
	if len(id) == 0 {
		return ErrInvalidID()
	}
	if err := h.taskService.DeleteTask(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return ErrResourceNotFound(err.Error())
		}
		return err
	}
	return c.JSON(fiber.Map{"deleted": id})
}

func (h *TaskHandler) HandleAssignTaskToSelf(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	user, err := getUserAuth(c)
	if err != nil {
		return err
	}
	params := types.UpdateTaskRequest{
		TaskID: id,
		UserID: user.ID,
	}
	if err := h.taskService.AssignTaskToSelf(c.Context(), params); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return ErrResourceNotFound(err.Error())
		}
		return err
	}
	return c.JSON(fiber.Map{"assigned": "true"})
}
func (h *TaskHandler) HandleAssignTaskToUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	var req types.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrBadRequest()
	}
	req.TaskID = id
	if err := h.taskService.AssignTaskToUser(c.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return ErrResourceNotFound(err.Error())
		case errors.Is(err, service.ErrTaskNotFound):
			return ErrResourceNotFound(err.Error())
		default:
			return err
		}
	}
	return c.JSON(fiber.Map{"assigned": "true"})
}
func (h *TaskHandler) HandlePutDueDateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return ErrInvalidID()
	}
	var req types.UpdateDueDateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrBadRequest()
	}
	if err := req.Validate(); err != nil {
		return ErrBadRequestCustomMessage(err.Error())
	}
	user, err := getUserAuth(c)
	if err != nil {
		return err
	}
	req.AssignedTo = user.ID
	if err := h.taskService.UpdateDueDate(c.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, service.ErrUnAuthorized):
			return ErrUnAuthorized()
		case errors.Is(err, service.ErrTaskNotFound):
			return ErrResourceNotFound(err.Error())
		default:
			return err
		}
	}
	return c.JSON(fiber.Map{"updated": id})
}
