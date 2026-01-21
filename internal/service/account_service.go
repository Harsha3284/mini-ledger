package service

import (
	"context"
	"errors"
	"mini-ledger/internal/model"
	"mini-ledger/internal/store"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid input")

type AccountService struct {
	store *store.AccountStore
}

func NewAccountService(s *store.AccountStore) *AccountService {
	return &AccountService{store: s}
}

func (s *AccountService) Create(ctx context.Context, name string, typ model.AccountType, cur model.Currency) (model.Account, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 80 {
		return model.Account{}, ErrInvalidInput
	}

	if typ != model.AccountChecking && typ != model.AccountSavings && typ != model.AccountCredit {
		return model.Account{}, ErrInvalidInput
	}
	if cur != model.USD {
		return model.Account{}, ErrInvalidInput
	}

	a := model.Account{
		ID:       uuid.NewString(),
		Name:     name,
		Type:     typ,
		Currency: cur,
	}
	return s.store.Create(ctx, a)
}

func (s *AccountService) Get(ctx context.Context, id string) (model.Account, error) {
	return s.store.GetByID(ctx, id)
}

func (s *AccountService) List(ctx context.Context, limit int) ([]model.Account, error) {
	return s.store.List(ctx, limit)
}

func (s *AccountService) UpdateName(ctx context.Context, id string, name string) (model.Account, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 80 {
		return model.Account{}, ErrInvalidInput
	}
	return s.store.UpdateName(ctx, id, name)
}

func (s *AccountService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
