package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mini-ledger/internal/model"
	"mini-ledger/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateEntryRequest struct {
	Direction   model.EntryDirection `json:"direction" binding:"required"`
	Amount      string               `json:"amount" binding:"required"`
	Category    *string              `json:"category"`
	Description *string              `json:"description"`
	OccurredAt  *time.Time           `json:"occurred_at"`
}

type LedgerHandler struct {
	svc *service.LedgerService
}

func NewLedgerHandler(svc *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{svc: svc}
}

func (h *LedgerHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/accounts/:id/entries", h.createEntry)
	rg.GET("/accounts/:id/entries", h.listEntries)
	rg.GET("/accounts/:id/balance", h.getBalance)
}

func (h *LedgerHandler) createEntry(c *gin.Context) {
	accountID := c.Param("id")

	var req CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	occurredAt := time.Time{}
	if req.OccurredAt != nil {
		occurredAt = *req.OccurredAt
	}

	entry, err := h.svc.CreateEntry(
		c.Request.Context(),
		accountID,
		req.Direction,
		req.Amount,
		req.Category,
		req.Description,
		occurredAt,
	)
	if err != nil {
		if err == service.ErrInvalidEntry {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry"})
			return
		}
		if err == service.ErrAccountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *LedgerHandler) listEntries(c *gin.Context) {
	accountID := c.Param("id")

	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	items, err := h.svc.ListEntries(c.Request.Context(), accountID, limit)
	if err != nil {
		if err == service.ErrAccountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *LedgerHandler) getBalance(c *gin.Context) {
	accountID := c.Param("id")

	bal, err := h.svc.GetBalance(c.Request.Context(), accountID)
	if err != nil {
		if err == service.ErrAccountNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account_id": accountID, "balance": bal})
}
