package api

import (
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
	//TODO: review
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
