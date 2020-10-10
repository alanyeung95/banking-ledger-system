package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/alanyeung95/banking-ledger-system/pkg/transactions"
)

// NewHandler return handler that serves the account service
func NewHandler(srv Service) http.Handler {
	h := handlers{srv}
	r := chi.NewRouter()
	//r.Get("/{id}", h.handleGetAccountBalance)
	r.Post("/create", h.handleCreateAccount)
	r.Post("/create-admin", h.handleCreateAdminAccount)
	r.Post("/transaction", h.handleTransaction)

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

func (h *handlers) handleTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var transaction transactions.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		kithttp.DefaultErrorEncoder(ctx,err, w)
		return
	}

	switch transaction.Operation {
	case transactions.Deposit:
		account, err := h.svc.GetAccountByID(ctx, r, transaction.Body.To)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}
		account, err = h.svc.UpdateBalance(ctx, r, transaction.Body.To, transaction.Body.Amount )
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}
		kithttp.EncodeJSONResponse(ctx, w, account)
		return
	case transactions.Withdraw:
		account, err := h.svc.GetAccountByID(ctx, r, transaction.Body.From)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}
		account, err = h.svc.UpdateBalance(ctx, r, transaction.Body.From, -transaction.Body.Amount )
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}
		kithttp.EncodeJSONResponse(ctx, w, account)
		return
	case transactions.Transfer:
		_, err := h.svc.GetAccountByID(ctx, r, transaction.Body.From)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}	
		_, err = h.svc.GetAccountByID(ctx, r, transaction.Body.To)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}		
		account, err := h.svc.TransferBalance(ctx, r, transaction.Body.From, transaction.Body.To, transaction.Body.Amount )
		kithttp.EncodeJSONResponse(ctx, w, account)
	default:
		kithttp.EncodeJSONResponse(ctx, w, "error: unsupport operation")
		return 
	}


	//kithttp.EncodeJSONResponse(ctx, w, "123")
}
