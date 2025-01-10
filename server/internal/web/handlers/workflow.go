package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	"github.com/shivamsanju/uploader/internal/web/utils"
	g "github.com/shivamsanju/uploader/pkg/globals"
)

type workflowHandler struct {
	wfRepo repo.WorkflowRepo
	dsRepo repo.DataStoreRepo
}

func NewWorkflowHandler() *workflowHandler {
	return &workflowHandler{
		wfRepo: repo.NewWorkflowRepo(),
		dsRepo: repo.NewDataStoreRepo(),
	}
}

func (h *workflowHandler) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	cbs, err := h.wfRepo.GetWorkflows(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if cbs == nil {
		cbs = make([]models.Workflow, 0)
	}
	render.JSON(w, r, cbs)
}

func (h *workflowHandler) GetWorkflowByID(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "id")
	cb, err := h.wfRepo.GetWorkflow(r.Context(), workflowID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func (h *workflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating workflow")
	workflow := &models.Workflow{}
	if err := render.DecodeJSON(r.Body, workflow); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	workflow.CreatedBy = r.Header.Get("email")
	workflow.UpdatedBy = r.Header.Get("email")

	g.Log.Infof("adding Workflow: %+v", workflow)
	id, err := h.wfRepo.CreateWorkflow(r.Context(), workflow)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (h *workflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "id")
	h.wfRepo.DeleteWorkflow(r.Context(), workflowID)
}
