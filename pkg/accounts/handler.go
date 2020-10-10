package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	kithttp "github.com/go-kit/kit/transport/http"
)

// NewHandler return handler that serves the account service
func NewHandler(srv Service) http.Handler {
	h := handlers{srv}
	r := chi.NewRouter()
	r.Get("/{id}", h.handleGetAccountBalance)
	r.Post("/create", h.handleCreateAccount)
	r.Post("/create-admin", h.handleCreateAdminAccount)

	return r
}

type handlers struct {
	svc Service
}

func (h *handlers) handleGetAccountBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//id := chi.URLParam(r, "id")

	kithttp.EncodeJSONResponse(ctx, w, "wip")
}

func (h *handlers) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var model Account
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		kithttp.DefaultErrorEncoder(ctx,err, w)
		return
	}

	response, err := h.svc.CreateAccount(ctx, r, &model)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, response)
}

func (h *handlers) handleCreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var model Account
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		kithttp.DefaultErrorEncoder(ctx,err, w)
		return
	}

	response, err := h.svc.CreateAdminAccount(ctx, r, &model)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, response)
}
