package accounts

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/alanyeung95/banking-ledger-system/pkg/transactions"
)

// NewHandler return handler that serves the account service
func NewHandler(srv Service, transSrv transactions.Service) http.Handler {
	h := handlers{srv, transSrv}
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Crypto.com Ops Team Back End Engineering Coding Cahllenge"))
	})

	r.Get("/accounts/{id}", h.handleGetAccountBalance)
	r.Post("/accounts/create", h.handleCreateAccount)
	r.Post("/accounts/create-admin", h.handleCreateAdminAccount)

	r.Post("/transactions", h.handleTransaction)
	r.Post("/transactions/{id}/undo", h.handleUndoTransaction)
	r.Get("/transactions", h.handleGetTransactionsByID)	
	return r
}

type handlers struct {
	svc Service
	transSrv transactions.Service
}

func (h *handlers) handleGetAccountBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	account, err := h.svc.GetAccountByID(ctx, r, id)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(account))
}

func (h *handlers) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var model Account
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		kithttp.DefaultErrorEncoder(ctx,err, w)
		return
	}

	newAccount, err := h.svc.CreateAccount(ctx, r, &model)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(newAccount))
}

func (h *handlers) handleCreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var model Account
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		kithttp.DefaultErrorEncoder(ctx,err, w)
		return
	}

	newAccount, err := h.svc.CreateAdminAccount(ctx, r, &model)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(newAccount))
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
		kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(account))

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
		kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(account))

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
		kithttp.EncodeJSONResponse(ctx, w, convertToAccountReadModel(account))

	default:
		kithttp.EncodeJSONResponse(ctx, w, "error: unsupport operation")
		return 
	}

	transaction.Time = time.Now()
	_, err := h.transSrv.RecordTransaction(ctx, r, &transaction)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
}

func (h *handlers) handleUndoTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	transaction, err := h.transSrv.GetTransactionByID(ctx, r, id)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}

	switch transaction.Operation {
	case transactions.Deposit:
		account, err := h.svc.UpdateBalance(ctx, r, transaction.Body.To, - transaction.Body.Amount )
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}

		transaction.Operation = transactions.Withdraw
		transaction.Body.From = transaction.Body.To
		transaction.Body.To = ""

		transaction.Time = time.Now()
		_, err = h.transSrv.RecordTransaction(ctx, r, transaction)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}

		kithttp.EncodeJSONResponse(ctx, w, account)		

	case transactions.Withdraw:
		account, err := h.svc.UpdateBalance(ctx, r, transaction.Body.From, transaction.Body.Amount )
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}

		transaction.Operation = transactions.Deposit
		transaction.Body.To = transaction.Body.From
		transaction.Body.From = ""

		transaction.Time = time.Now()
		_, err = h.transSrv.RecordTransaction(ctx, r, transaction)
		if err != nil {
			kithttp.DefaultErrorEncoder(ctx, err, w)
			return
		}

		kithttp.EncodeJSONResponse(ctx, w, account)	

	default:
		kithttp.EncodeJSONResponse(ctx, w, "error: unsupport operation, only withdraw or deposit operation is allowed to be undo")
		return 		
	}

	transaction.Operation = transactions.Undo
	transaction.Body.From = ""
	transaction.Body.To = ""
	transaction.Body.Amount = 0
	transaction.Notes = id

	_, err = h.transSrv.RecordTransaction(ctx, r, transaction)
	if err != nil {
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
}

func (h *handlers) handleGetTransactionsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
    asc, err := strconv.Atoi( r.URL.Query().Get("asc"))
    if err != nil  {
		asc = -1
	}	

	transactionList, err := h.transSrv.GetTransactionsByID(ctx, r, id, asc)
	if err != nil{
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, transactionList)
}

func convertToAccountReadModel(account *Account) *AccountReadModel {
	return &AccountReadModel{
		ID: account.ID,
		Name: account.Name,
		Balance: account.Balance,
	}
}
