package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ismail118/simple-bank/repository"
	"github.com/ismail118/simple-bank/token"
	"github.com/ismail118/simple-bank/util"
)

// Server serves HTTP request for our banking service
type Server struct {
	store      repository.Store
	repo       repository.Repository
	router     *gin.Engine
	tokenMaker token.Maker
	config     *util.Config
}

// NewServer create a new HTTP server and setup routing
func NewServer(
	store repository.Store,
	repo repository.Repository,
	tokenMaker token.Maker,
	config *util.Config,
) Server {
	server := Server{
		store:      store,
		repo:       repo,
		router:     nil,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (server *Server) setupRouter() {
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

	//auth := router.Group("/api")
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUsers)
	router.GET("/users", server.listUsers)
	router.PUT("/users", server.updateUsers)
	router.DELETE("/users/:username", server.deleteUsers)

	router.POST("/users/login", server.loginUser)

	server.router = router
}
