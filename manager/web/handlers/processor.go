package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	commonutils "github.com/uploadpilot/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc/processor"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type processorHandler struct {
	pSvc *processor.Service
}

func NewProcessorsHandler(pSvc *processor.Service) *processorHandler {
	return &processorHandler{
		pSvc: pSvc,
	}
}

func (h *processorHandler) GetProcessors(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	processors, err := h.pSvc.GetAllProcessorsInWorkspace(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, processors)
}

func (h *processorHandler) GetProcessorDetailsByID(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	processor, err := h.pSvc.GetProcessor(r.Context(), processorID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, processor)
}
func (h *processorHandler) GetTemplates(w http.ResponseWriter, r *http.Request) ([]dto.ProcessorTemplate, error) {
	templates := h.pSvc.GetTemplates(r.Context())
	return templates, nil
}

func (h *processorHandler) CreateProcessor(w http.ResponseWriter, r *http.Request, body *dto.CreateProcessorRequest) (*string, error) {
	workspaceID := chi.URLParam(r, "workspaceId")

	var processor models.Processor
	if err := commonutils.ConvertDTOToModel(body, &processor); err != nil {
		return nil, err
	}
	if err := h.pSvc.CreateProcessor(r.Context(), workspaceID, &processor, &body.TemplateKey); err != nil {
		return nil, err
	}

	return &processor.ID, nil
}

func (h *processorHandler) UpdateProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	workspaceID := chi.URLParam(r, "workspaceId")

	processor := &dto.EditProcRequest{}
	if err := render.DecodeJSON(r.Body, processor); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.pSvc.EditNameAndTrigger(r.Context(), workspaceID, processorID, processor)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) DeleteProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	workspaceID := chi.URLParam(r, "workspaceId")

	if err := h.pSvc.DeleteProcessor(r.Context(), workspaceID, processorID); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	workspaceID := chi.URLParam(r, "workspaceId")

	workflow := &dto.WorkflowUpdate{}
	if err := render.DecodeJSON(r.Body, workflow); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	if err := h.pSvc.UpdateWorkflow(r.Context(), workspaceID, processorID, workflow.Workflow); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) EnableProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	workspaceID := chi.URLParam(r, "workspaceId")

	err := h.pSvc.EnableDisableProcessor(r.Context(), workspaceID, processorID, true)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) DisableProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	workspaceID := chi.URLParam(r, "workspaceId")

	err := h.pSvc.EnableDisableProcessor(r.Context(), workspaceID, processorID, false)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, h.pSvc.GetAllTasks(r.Context()))
}

func (h *processorHandler) GetWorkflowRuns(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	processorID := chi.URLParam(r, "processorId")
	return h.pSvc.GetWorkflowRuns(r.Context(), processorID)
}

func (h *processorHandler) GetWorkflowLogs(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	workflowID := r.URL.Query().Get("workflowId")
	runID := r.URL.Query().Get("runId")
	return h.pSvc.GetWorkflowHistory(r.Context(), workflowID, runID)
}
