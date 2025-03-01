package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/processor"
	"github.com/uploadpilot/manager/internal/utils"
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
func (h *processorHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	templates := h.pSvc.GetTemplates(r.Context())
	render.JSON(w, r, templates)
}

func (h *processorHandler) CreateProcessor(r *http.Request, params dto.WorkspaceParams, query interface{}, body dto.CreateProcessorRequest) (*string, int, error) {
	var processor models.Processor
	if err := utils.ConvertDTOToModel(&body, &processor); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}
	if err := h.pSvc.CreateProcessor(r.Context(), params.WorkspaceID, &processor, &body.TemplateKey); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &processor.ID, http.StatusOK, nil
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

func (h *processorHandler) GetWorkflowRuns(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")
	runs, err := h.pSvc.GetWorkflowRuns(r.Context(), processorID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, runs)
}

func (h *processorHandler) GetWorkflowLogs(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Query().Get("workflowId")
	runID := r.URL.Query().Get("runId")
	details, err := h.pSvc.GetWorkflowHistory(r.Context(), workflowID, runID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, details)
}
