package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/tasks"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type processorHandler struct {
	pSvc *svc.ProcessorService
}

func NewProcessorsHandler(pSvc *svc.ProcessorService) *processorHandler {
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

func (h *processorHandler) CreateProcessor(w http.ResponseWriter, r *http.Request, body *dto.CreateProcessorRequest) (*string, error) {
	workspaceID := chi.URLParam(r, "workspaceId")

	var processor models.Processor
	if err := commonutils.ConvertDTOToModel(body, &processor); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
	}
	if err := h.pSvc.CreateProcessor(r.Context(), workspaceID, &processor); err != nil {
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

func (h *processorHandler) GetProcBlock(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, tasks.GetAllTasks())
}
