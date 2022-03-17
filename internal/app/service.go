package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/strick-j/scimfe/internal/config"
	"github.com/strick-j/scimfe/internal/repository"
	"github.com/strick-j/scimfe/internal/service"
	"github.com/strick-j/scimfe/internal/web"
	"github.com/strick-j/scimfe/internal/web/handler"
	"github.com/strick-j/scimfe/internal/web/middleware"
	"go.uber.org/zap"
)

type Service struct {
	server *web.Server
	logger *zap.Logger
}

func NewService(baseCtx context.Context, logger *zap.Logger, conn *Connectors, cfg *config.Config) *Service {
	srv := web.NewServer(cfg.Server.ListenParams())

	userStore := repository.NewUserRepository(conn.DB)
	sessionStore := repository.NewSessionRepository(conn.Redis)

	userSvc := service.NewUsersService(logger, userStore)
	authSvc := service.NewAuthService(logger, userSvc, sessionStore)

	hWrapper := web.NewWrapper(logger.Named("http"))
	requireAuth := hWrapper.MiddlewareFunc(middleware.NewAuthMiddleware(authSvc))

	// General
	srv.Router.Methods(http.MethodGet).
		Path("/ping").
		HandlerFunc(hWrapper.WrapResourceHandler(handler.Ping))

	// Auth
	authHandler := handler.NewAuthHandler(userSvc, authSvc)
	srv.Router.Methods(http.MethodPost).
		Path("/auth").
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.Login))
	srv.Router.Methods(http.MethodPost).
		Path("/auth/register").
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.Register))

	// Session
	sessionRouter := srv.Router.Path("/auth/session").Subrouter()
	sessionRouter.Use(requireAuth)
	sessionRouter.Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(authHandler.GetSession))
	sessionRouter.Methods(http.MethodDelete).
		HandlerFunc(hWrapper.WrapHandler(authHandler.Logout))

	// Users
	usrHandler := handler.NewUserHandler(userSvc)
	usrRouter := srv.Router.NewRoute().Subrouter()
	usrRouter.Use(requireAuth)
	usrRouter.Path("/users").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetUsersList))
	usrRouter.Path("/users/self").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetCurrentUser))
	//usrRouter.Path("/users/self/balance").Methods(http.MethodGet).
	//	HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetBalance))
	usrRouter.Path("/users/{userId}").Methods(http.MethodGet).
		HandlerFunc(hWrapper.WrapResourceHandler(usrHandler.GetByID))

	return &Service{
		server: srv,
		logger: logger,
	}
}

// Start starts the service
func (s Service) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Info("starting http service", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(ctx); err != nil {
			if err == context.Canceled {
				return
			}
			s.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	wg.Wait()
	s.logger.Info("goodbye")
}
