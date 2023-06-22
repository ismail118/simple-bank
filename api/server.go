package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/repository"
)

// Server serves HTTP request for our banking service
type Server struct {
	store  repository.Store
	repo   repository.Repository
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing
func NewServer(store repository.Store, repo repository.Repository) Server {
	server := Server{
		store:  store,
		repo:   repo,
		router: nil,
	}

	router := gin.Default()

	// add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.GET("/entries/:id", server.getEntry)
	router.GET("/entries", server.listEntries)

	router.GET("/transfer/:id", server.getTransfer)
	router.GET("/transfer", server.listTransfer)
	router.POST("/transfer", server.transfer)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
