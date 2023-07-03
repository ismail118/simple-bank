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
	// no need auth
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	// need auth
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/accounts", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.GET("/entries/:id", server.getEntry)
	authRoutes.GET("/entries", server.listEntries)

	authRoutes.GET("/transfer/:id", server.getTransfer)
	authRoutes.GET("/transfer", server.listTransfer)
	authRoutes.POST("/transfer", server.transfer)

	authRoutes.GET("/users/:username", server.getUsers)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.PUT("/users", server.updateUsers)
	authRoutes.DELETE("/users/:username", server.deleteUsers)

	server.router = router
}
