package handlers

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/uploadpilot/core/internal/activities/catalog"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
)

type processorHandler struct {
	pSvc *services.ProcessorService
}

func NewProcessorsHandler(pSvc *services.ProcessorService) *processorHandler {
	return &processorHandler{
		pSvc: pSvc,
	}
}

func (h *processorHandler) GetProcessors(r *http.Request, params dto.WorkspaceParams, query, body interface{}) ([]models.Processor, int, error) {
	processors, err := h.pSvc.GetAllProcessorsInWorkspace(r.Context(), params.WorkspaceID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return processors, http.StatusOK, nil
}

func (h *processorHandler) GetProcessorDetailsByID(r *http.Request, params dto.ProcessorParams, query, body interface{}) (*models.Processor, int, error) {
	processor, err := h.pSvc.GetProcessor(r.Context(), params.ProcessorID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return processor, http.StatusOK, nil
}

func (h *processorHandler) GetTemplates(r *http.Request, params dto.WorkspaceParams, query, body interface{}) ([]dto.ProcessorTemplate, int, error) {
	templates := h.pSvc.GetTemplates(r.Context())
	return templates, http.StatusOK, nil
}

func (h *processorHandler) CreateProcessor(r *http.Request, params dto.WorkspaceParams, query interface{}, body dto.CreateProcessorRequest) (*string, int, error) {
	var processor models.Processor
	if err := copier.Copy(&processor, &body); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}
	if err := h.pSvc.CreateProcessor(r.Context(), params.WorkspaceID, &processor, &body.TemplateKey); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &processor.ID, http.StatusOK, nil
}

func (h *processorHandler) UpdateProcessor(r *http.Request, params dto.ProcessorParams, query interface{}, body dto.EditProcRequest) (*string, int, error) {
	err := h.pSvc.EditNameAndTrigger(r.Context(), params.WorkspaceID, params.ProcessorID, &body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return nil, http.StatusOK, nil
}

func (h *processorHandler) DeleteProcessor(r *http.Request, params dto.ProcessorParams, query, body interface{}) (bool, int, error) {
	if err := h.pSvc.DeleteProcessor(r.Context(), params.WorkspaceID, params.ProcessorID); err != nil {
		return false, http.StatusBadRequest, err
	}
	return true, http.StatusOK, nil
}

func (h *processorHandler) UpdateWorkflow(r *http.Request, params dto.ProcessorParams, query interface{}, body dto.WorkflowUpdate) (bool, int, error) {
	if err := h.pSvc.UpdateWorkflow(r.Context(), params.WorkspaceID, params.ProcessorID, body.Workflow); err != nil {
		return false, http.StatusBadRequest, err
	}
	return true, http.StatusOK, nil
}

func (h *processorHandler) EnableProcessor(r *http.Request, params dto.ProcessorParams, query, body interface{}) (bool, int, error) {
	err := h.pSvc.EnableDisableProcessor(r.Context(), params.WorkspaceID, params.ProcessorID, true)
	if err != nil {
		return false, http.StatusBadRequest, err
	}

	return true, http.StatusOK, nil
}

func (h *processorHandler) DisableProcessor(r *http.Request, params dto.ProcessorParams, query, body interface{}) (bool, int, error) {
	err := h.pSvc.EnableDisableProcessor(r.Context(), params.WorkspaceID, params.ProcessorID, false)
	if err != nil {
		return false, http.StatusBadRequest, err
	}

	return true, http.StatusOK, nil
}

func (h *processorHandler) GetAllActivities(r *http.Request, params dto.ProcessorParams, query, body interface{}) ([]catalog.ActivityMetadata, int, error) {
	return h.pSvc.GetAllActivities(r.Context()), http.StatusOK, nil
}

func (h *processorHandler) GetWorkflowRuns(r *http.Request, params dto.ProcessorParams, query, body interface{}) ([]dto.WorkflowRun, int, error) {
	runs, err := h.pSvc.GetWorkflowRuns(r.Context(), params.ProcessorID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return runs, http.StatusOK, nil
}

func (h *processorHandler) GetWorkflowLogs(r *http.Request, params dto.WorkflowRunParams, query, body interface{}) ([]dto.WorkflowRunLogs, int, error) {
	details, err := h.pSvc.GetWorkflowHistory(r.Context(), params.WorkflowID, params.RunID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return details, http.StatusOK, nil
}
