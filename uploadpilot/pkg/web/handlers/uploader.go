package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	"github.com/uploadpilot/uploadpilot/pkg/db/repo"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/web/utils"
)

type uploaderHandler struct {
	wfRepo repo.UploaderRepo
	dsRepo repo.DataStoreRepo
}

func NewuploaderHandler() *uploaderHandler {
	return &uploaderHandler{
		wfRepo: repo.NewUploaderRepo(),
		dsRepo: repo.NewDataStoreRepo(),
	}
}

func (h *uploaderHandler) GetAllUploaders(w http.ResponseWriter, r *http.Request) {
	cbs, err := h.wfRepo.GetAll(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if cbs == nil {
		cbs = make([]models.Uploader, 0)
	}
	render.JSON(w, r, cbs)
}

func (h *uploaderHandler) GetUploaderByID(w http.ResponseWriter, r *http.Request) {
	uploaderID := chi.URLParam(r, "uploaderId")
	cb, err := h.wfRepo.Get(r.Context(), uploaderID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func (h *uploaderHandler) CreateUploader(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating uploader")
	uploader := &models.Uploader{}
	if err := render.DecodeJSON(r.Body, uploader); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	uploader.CreatedBy = r.Header.Get("email")
	uploader.UpdatedBy = r.Header.Get("email")

	g.Log.Infof("adding uploader: %+v", uploader)
	id, err := h.wfRepo.Create(r.Context(), uploader)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (h *uploaderHandler) DeleteUploader(w http.ResponseWriter, r *http.Request) {
	uploaderID := chi.URLParam(r, "uploaderId")
	h.wfRepo.Delete(r.Context(), uploaderID)
}

func (h *uploaderHandler) GetAllAllowedSources(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []models.AllowedSources{
		models.FileUpload,
		models.Webcamera,
		models.Audio,
		models.ScreenCapture,
		models.Box,
		models.Dropbox,
		models.Facebook,
		models.GoogleDrive,
		models.GooglePhotos,
		models.Instagram,
		models.OneDrive,
		models.Unsplash,
		models.Url,
		models.Zoom,
	})
}

func (h *uploaderHandler) UpdateUploaderConfig(w http.ResponseWriter, r *http.Request) {
	uploaderID := chi.URLParam(r, "uploaderId")
	uploaderConfig := &models.UploaderConfig{}
	if err := render.DecodeJSON(r.Body, uploaderConfig); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	updatedBy := r.Header.Get("email")
	err := h.wfRepo.UpdateConfig(r.Context(), uploaderID, uploaderConfig, updatedBy)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, uploaderID)
}
