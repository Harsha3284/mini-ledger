package handlers

import "mini-ledger/internal/model"

type CreateAccountRequest struct {
	Name     string            `json:"name" binding:"required"`
	Type     model.AccountType `json:"type" binding:"required"`
	Currency model.Currency    `json:"currency" binding:"required"`
}

type UpdateAccountRequest struct {
	Name string `json:"name" binding:"required"`
}
