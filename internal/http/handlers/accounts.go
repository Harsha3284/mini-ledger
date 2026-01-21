package handlers

import (
	"net/http"
	"strconv"

	"mini-ledger/internal/service"
	"mini-ledger/internal/store"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/accounts", h.create)
	rg.GET("/accounts/:id", h.getByID)
	rg.GET("/accounts", h.list)
	rg.PATCH("/accounts/:id", h.update)
	rg.DELETE("/accounts/:id", h.delete)
}

func (h *AccountHandler) create(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	acc, err := h.svc.Create(c.Request.Context(), req.Name, req.Type, req.Currency)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (h *AccountHandler) getByID(c *gin.Context) {
	id := c.Param("id")

	acc, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *AccountHandler) list(c *gin.Context) {
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	items, err := h.svc.List(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *AccountHandler) update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	acc, err := h.svc.UpdateName(c.Request.Context(), id, req.Name)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *AccountHandler) delete(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		if err == store.ErrConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "account has ledger entries; cannot delete"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
