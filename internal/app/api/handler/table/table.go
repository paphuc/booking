package tablehandler

import (
	"context"
	"encoding/json"
	"net/http"

	"booking/configs"
	"booking/internal/app/types"
	"booking/internal/pkg/glog"
	"booking/internal/pkg/respond"

	"github.com/go-playground/validator/v10"
	//"github.com/gorilla/mux"
)

type (
	service interface {
		// Get(ctx context.Context, id string) (*types.Member, error)
		InsertTable(ctx context.Context, tableRequest types.TableRequest) (*types.Table , error)
		UpdateTableByID(ctx context.Context, Table types.UpdateTableRequest) error
		DeleteTable(ctx context.Context, Table types.DeleteTableRequest) error
	}

	// Handler is table web handler
	Handler struct {
		conf   *configs.Configs
		em     *configs.ErrorMessage
		srv    service
		logger glog.Logger
	}
)

var (
	validate = validator.New()
)

// New return new rest api table handler
func New(c *configs.Configs, e *configs.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

// Post hanlder insert table HTTP request
func (h *Handler) InsertTable(w http.ResponseWriter, r *http.Request) {

	var tableRequest types.TableRequest

	if err := json.NewDecoder(r.Body).Decode(&tableRequest); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(tableRequest); err != nil {
		h.logger.Errorf("Failed when validate field tableRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	tableReq, err := h.srv.InsertTable(r.Context(), tableRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, tableReq)
}

// Put hanlder update table HTTP request
func (h *Handler) UpdateTableByID(w http.ResponseWriter, r *http.Request) {

	var table types.UpdateTableRequest

	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateTableRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	if err := validate.Struct(table); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateTableRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if error := h.srv.UpdateTableByID(r.Context(), table); error != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Put hanlder delete table HTTP request
func (h *Handler) DeleteTable(w http.ResponseWriter, r *http.Request) {

	var table types.DeleteTableRequest

	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		h.logger.Errorf("Failed when validate field in method DeleteTableRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	if err := validate.Struct(table); err != nil {
		h.logger.Errorf("Failed when validate field in method DeleteTableRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if error := h.srv.DeleteTable(r.Context(), table); error != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}


