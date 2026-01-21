package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not-ready", "error": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	return r
}