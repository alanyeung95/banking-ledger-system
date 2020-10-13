package accounts

import (
	"context"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service interface
type Service interface {
	CreateAccount(ctx context.Context, r *http.Request, user *Account) (*Account, error)
	CreateAdminAccount(ctx context.Context, r *http.Request, user *Account) (*Account, error)
	GetAccountByID(ctx context.Context, r *http.Request, id string) (*Account, error)
	UpdateBalance(ctx context.Context, r *http.Request, id string, amount int) (*Account, error)
	TransferBalance(ctx context.Context, r *http.Request, from string, to string, amount int) (*Account, error)
	DeleteAccountByID(ctx context.Context, r *http.Request, from string) error
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

func (s *service) GetAccountByID(ctx context.Context, r *http.Request, id string) (*Account, error) {
	return s.repository.Find(ctx, id)
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

func (s *service) UpdateBalance(ctx context.Context, r *http.Request, id string, amount int) (*Account, error) {
	return s.repository.UpdateBalance(ctx, id, amount)
}

func (s *service) TransferBalance(ctx context.Context, r *http.Request, from string, to string, amount int) (*Account, error) {
	sourceAccount, err := s.UpdateBalance(ctx, r, from, -amount)
	if err != nil {
		return nil, err
	}
	_, err = s.UpdateBalance(ctx, r, to, amount)
	if err != nil {
		return nil, err
	}
	return sourceAccount, nil
}

func (s *service) DeleteAccountByID(ctx context.Context, r *http.Request, id string) error {
	return s.repository.Delete(ctx, id)
}
