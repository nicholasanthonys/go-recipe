package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/adapter/repository"
	"github.com/nicholasanthonys/go-recipe/adapter/validator"
)

type ginEngine struct {
	router     *gin.Engine
	log        logger.Logger
	db         repository.NoSQL
	kv         repository.KeyValStoreIn
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newGinServer(
	log logger.Logger,
	db repository.NoSQL,
	kv repository.KeyValStoreIn,
	validator validator.Validator,
	port Port,
	t time.Duration,
) *ginEngine {
	return &ginEngine{
		router:     gin.Default(),
		log:        log,
		db:         db,
		kv:         kv,
		validator:  validator,
		port:       port,
		ctxTimeout: t,
	}
}

func (g ginEngine) Listen() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()

	g.setAppHandlers(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
		if err := server.ListenAndServe(); err != nil {
			g.log.WithError(err).Fatalln("Error starting HTTP server")
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		g.log.WithError(err).Fatalln("Server Shutdown Failed")
	}

	g.log.Infof("Service down")
}

/* TODO ADD MIDDLEWARE */
func (g ginEngine) setAppHandlers(router *gin.Engine) {

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// auth, err := NewAuthenticator()
	// if err != nil {
	// 	g.log.Fatalln("Failed to initialize the authenticator: %v", err)
	// }

	// router.GET("/auth0/login", Auth0HandlerLogin(auth))
	// router.GET("/auth0/callback", Auth0HandlerCallback(auth))
	// router.GET("/auth0/user", Auth0IsAuthenticated, Auth0UserProfileHandler)
	// router.GET("/auth0/logout", Auth0HandlerLogout)

	router.POST("/v1/signin", g.SignInHandler)
	router.POST("/v1/signout", g.SignOutHandler)

	router.POST("/v1/register", g.RegisterHandler)

	router.POST("/v1/refresh", g.RefreshHandler)

	// router.Use(g.CheckAPIKey)
	// router.Use(g.AuthMiddlewareWithSession)

	// auth 0 middleware
	router.Use(Auth0Middleware)

	router.POST("/v1/transfers", g.buildCreateTransferAction())
	router.GET("/v1/transfers", g.buildFindAllTransferAction())

	router.GET("/v1/accounts/:account_id/balance", g.buildFindBalanceAccountAction())
	router.POST("/v1/accounts", g.buildCreateAccountAction())
	router.GET("/v1/accounts", g.buildFindAllAccountAction())

	router.POST("/v1/recipes", g.NewRecipeHandler)
	router.GET("/v1/recipes", g.ListRecipesHandler)
	router.GET("/v1/recipes/:id", g.FindRecipeHandler)

	router.PUT("/v1/recipes/:id", g.UpdateRecipeHandler)
	router.DELETE("/v1/recipes/:id", g.DeleteRecipeHandler)
	router.GET("/v1/recipes/search", g.SearchRecipeHandler)

	router.GET("/v1/health", g.healthcheck())
}
