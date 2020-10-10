package accounts

import (
	"fmt"
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service interface
type Service interface {
	CreateAccount(ctx context.Context, r *http.Request, user *Account) (*Account, error)
	CreateAdminAccount(ctx context.Context, r *http.Request, user *Account) (*Account, error)
}

type service struct {
	repository Repository
}

// NewService start the new service
func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateAccount(ctx context.Context, r *http.Request, account *Account) (*Account, error) {
	account.ID = uuid.NewV4().String()
	account.Password = hashAndSalt([]byte(account.Password))
	account.UserGroup = GeneralUser
	return s.repository.Upsert(ctx, account.ID, *account)
}

func (s *service) CreateAdminAccount(ctx context.Context, r *http.Request, account *Account) (*Account, error) {
	account.ID = uuid.NewV4().String()
	account.Password = hashAndSalt([]byte(account.Password))
	account.UserGroup = OperationTeam
	return s.repository.Upsert(ctx, account.ID, *account)
}

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		// todo: error handling
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}