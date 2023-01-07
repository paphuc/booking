package memberhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"booking/configs"
	"booking/internal/app/types"
	"booking/internal/pkg/glog"
	"booking/internal/pkg/respond"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	service interface {
		Get(ctx context.Context, id string) (*types.Member, error)
		InsertMember(ctx context.Context, MemberRequest types.MemberRequest) (*types.Member, error)
		UpdateMemberByID(ctx context.Context, Member types.UpdateMemberRequest) error
		Login(ctx context.Context, MemberLogin types.MemberLogin) (*types.MemberResponseSignUp, error)
	}

	// Handler is member web handler
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

// New return new rest api member handler
func New(c *configs.Configs, e *configs.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

// Get handle get member HTTP request
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	member, err := h.srv.Get(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, member)
}

// Post hanlder insert member HTTP request
func (h *Handler) InsertMember(w http.ResponseWriter, r *http.Request) {

	var memberRequest types.MemberRequest

	if err := json.NewDecoder(r.Body).Decode(&memberRequest); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(memberRequest); err != nil {
		h.logger.Errorf("Failed when validate field memberRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	mem, err := h.srv.InsertMember(r.Context(), memberRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, mem)
}

// Put hanlder update member HTTP request
func (h *Handler) UpdateMemberByID(w http.ResponseWriter, r *http.Request) {

	var member types.UpdateMemberRequest

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateMemberByID", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	if err := validate.Struct(member); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateMemberByID", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if error := h.srv.UpdateMemberByID(r.Context(), member); error != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var MemberLogin types.MemberLogin

	if err := json.NewDecoder(r.Body).Decode(&MemberLogin); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(MemberLogin); err != nil {
		h.logger.Errorf("Failed when validate field MemberLogin", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	member, err := h.srv.Login(r.Context(), MemberLogin)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.IncorrectPasswordEmail)
		return
	}

	respond.JSON(w, http.StatusOK, member)
}
