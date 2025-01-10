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

type importPolicyHandler struct {
	ipRepo repo.ImportPolicyRepo
}

func NewImportPolicyHandler() importPolicyHandler {
	return importPolicyHandler{
		ipRepo: repo.NewImportPolicyRepo(),
	}
}

func (ip *importPolicyHandler) GetImportPolicies(w http.ResponseWriter, r *http.Request) {
	cbs, err := ip.ipRepo.GetImportPolicies(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if cbs == nil {
		cbs = make([]models.ImportPolicy, 0)
	}
	render.JSON(w, r, cbs)
}

func (ip *importPolicyHandler) GetImportPolicy(w http.ResponseWriter, r *http.Request) {
	importPolicyID := chi.URLParam(r, "id")
	cb, err := ip.ipRepo.GetImportPolicy(r.Context(), importPolicyID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func (ip *importPolicyHandler) CreateImportPolicy(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating import policy")
	importPolicy := &models.ImportPolicy{AllowedMimeTypes: []string{"image/png", "image/jpeg"}}
	if err := render.DecodeJSON(r.Body, importPolicy); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	g.Log.Infof("adding import policy: %+v", importPolicy)
	importPolicy.CreatedBy = r.Header.Get("email")
	importPolicy.UpdatedBy = r.Header.Get("email")
	id, err := ip.ipRepo.CreateImportPolicy(r.Context(), importPolicy)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (ip *importPolicyHandler) DeleteImportPolicy(w http.ResponseWriter, r *http.Request) {
	importPolicyID := chi.URLParam(r, "id")
	ip.ipRepo.DeleteImportPolicy(r.Context(), importPolicyID)
}
