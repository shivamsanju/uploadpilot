package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/proc"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type processorHandler struct {
	processorSvc *proc.ProcessorService
}

func NewProcessorsHandler() *processorHandler {
	return &processorHandler{
		processorSvc: proc.NewProcessorService(),
	}
}

func (h *processorHandler) GetProcessors(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	processors, err := h.processorSvc.GetAllProcessorsInWorkspace(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, processors)
}

func (h *processorHandler) GetProcessorDetailsByID(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	processor, err := h.processorSvc.GetProcessor(r.Context(), processorID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	// p := struct {
	// 	*models.Processor
	// 	time int64
	// }{
	// 	Processor: processor,
	// 	time:      time.Now().Unix(),
	// }

	render.JSON(w, r, processor)
}

func (h *processorHandler) CreateProcessor(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	processor := &models.Processor{}
	if err := render.DecodeJSON(r.Body, processor); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.processorSvc.CreateProcessor(r.Context(), workspaceID, processor); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, processor.ID)
}

func (h *processorHandler) UpdateProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	processor := &dto.EditProcRequest{}
	if err := render.DecodeJSON(r.Body, processor); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.processorSvc.EditNameAndTrigger(r.Context(), processorID, processor)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) DeleteProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	if err := h.processorSvc.DeleteProcessor(r.Context(), processorID); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) UpdateTasks(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	tasks := &models.ProcTaskCanvas{}
	if err := render.DecodeJSON(r.Body, tasks); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.processorSvc.UpdateTasks(r.Context(), processorID, tasks); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) EnableProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	err := h.processorSvc.EnableDisableProcessor(r.Context(), processorID, true)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) DisableProcessor(w http.ResponseWriter, r *http.Request) {
	processorID := chi.URLParam(r, "processorId")

	err := h.processorSvc.EnableDisableProcessor(r.Context(), processorID, false)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *processorHandler) GetProcBlock(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, proc.ProcTaskBlocks[1:])
}
