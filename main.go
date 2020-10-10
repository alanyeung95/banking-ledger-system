package main	

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/alanyeung95/banking-ledger-system/pkg/mongo"
	"github.com/alanyeung95/banking-ledger-system/pkg/accounts"
	"github.com/alanyeung95/banking-ledger-system/pkg/transactions"


)
func main() {
	// todo: maybe remove hardcoding
	mongoClient, err := mongo.NewClient(
		"localhost",
		"",
		"",
		"banking",
	)

	if err != nil {
		fmt.Println("Cannot connect to mongoDB: "+ err.Error())
		return
	}

	accountSrv, err := newAccountSrv(mongoClient)
	if err != nil {
		fmt.Println("Cannot initialize account service: " + err.Error())
		return
	}

	
	transSrv, err := newTransactionSrv(mongoClient)
	if err != nil {
		fmt.Println("Cannot initialize account service: " + err.Error())
		return
	}


	r := chi.NewRouter()

	// Route - API
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Crypto.com Ops Team Back End Engineering Coding Cahllenge"))
		})
		r.Mount("/accounts", accounts.NewHandler(accountSrv))
		r.Mount("/services", transactions.NewHandler(transSrv))
	})

	addr := fmt.Sprintf(":%d", 3000)
	http.ListenAndServe(addr, r)
}

func newAccountSrv(client *mongodriver.Client) (accounts.Service, error) {
	accountRepository, err := mongo.NewAccountRepository(client, "banking", "accounts")
	if err != nil {
		return nil, err
	}

	srv := accounts.NewService(accountRepository)

	return srv, nil
}

func newTransactionSrv(client *mongodriver.Client) (transactions.Service, error) {
	transactionRepository, err := mongo.NewTransactionRepository(client, "banking", "transactions")
	if err != nil {
		return nil, err
	}

	srv := transactions.NewService( transactionRepository)

	return srv, nil
}
 