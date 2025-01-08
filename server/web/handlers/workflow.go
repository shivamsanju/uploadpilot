package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"github.com/shivamsanju/uploader/web/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkflowHandler interface {
	CreateWorkflow(w http.ResponseWriter, r *http.Request)
	ListWorkflows(w http.ResponseWriter, r *http.Request)
	GetWorkflow(w http.ResponseWriter, r *http.Request)
}

type workflowHandler struct {
}

func NewWorkflowHandler() WorkflowHandler {
	return &workflowHandler{}
}

func (h *workflowHandler) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	cbs, err := repo.GetWorkflows(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if cbs == nil {
		cbs = make([]models.Workflow, 0)
	}
	render.JSON(w, r, cbs)
}

func (h *workflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	WorkflowId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	cb, err := repo.GetWorkflow(r.Context(), WorkflowId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func (h *workflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating workflow")
	body := &models.Workflow{}
	if err := render.DecodeJSON(r.Body, body); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	body.UpdatedAt = time.Now().UTC().Unix()
	body.UpdatedAt = time.Now().UTC().Unix()
	body.CreatedBy = "John Doe"
	body.UpdatedBy = "John Doe"
	cb := body

	g.Log.Infof("adding Workflow: %+v", cb)
	id := repo.AddWorkflow(r.Context(), cb)

	render.JSON(w, r, id)
}
