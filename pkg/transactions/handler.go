package transactions

import (
	//"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

// NewHandler return handler that serves the account service
func NewHandler(srv Service) http.Handler {
	h := handlers{srv}
	r := chi.NewRouter()
	r.Get("/", h.handleGetTransactionsByID)

	return r
}

type handlers struct {
	svc Service
}

func (h *handlers) handleGetTransactionsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
    asc, err := strconv.Atoi( r.URL.Query().Get("asc"))
    if err != nil  {
		asc = -1
	}	
	
	transactionList, err := h.svc.GetTransactionsByID(ctx, r, id, asc)
	if err != nil{
		kithttp.DefaultErrorEncoder(ctx, err, w)
		return
	}
	kithttp.EncodeJSONResponse(ctx, w, transactionList)
}
