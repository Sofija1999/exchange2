package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/config"
	"github.com/Bloxico/exchange-gateway/sofija/core/usecases"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/Bloxico/exchange-gateway/sofija/handlers/product"
	"github.com/Bloxico/exchange-gateway/sofija/handlers/user"
	"github.com/Bloxico/exchange-gateway/sofija/log"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/Bloxico/exchange-gateway/sofija/server/auth"
	restful "github.com/emicklei/go-restful/v3"
)

type Server struct {
	srv    *http.Server
	wsCont *restful.Container

	RequestLogger log.Logger
}

type ApiVersion int

const (
	V1 ApiVersion = iota
	V2
)

func NewServer(cfg config.ServerConfig, db *database.DB) *Server {

	// http Server
	httpSrv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// restful container, where all the services and routes will be connected
	wsCont := restful.NewContainer()

	// the server, encapsulating all services, logging, etc.
	fullSrv := &Server{
		srv:    httpSrv,
		wsCont: wsCont,

		RequestLogger: cfg.Logger,
	}

	// add logging
	wsCont.Filter(log.NCSACommonLogFormatLogger(cfg.Logger))

	// add error handling
	wsCont.DoNotRecover(false)
	wsCont.ServiceErrorHandler(fullSrv.WriteServiceErrorJson)
	wsCont.RecoverHandler(fullSrv.RecoverHandler)

	// base server paths
	baseWs := new(restful.WebService)
	baseWs.Path("/")
	baseWs.Route(baseWs.GET("/ping").Filter(auth.AuthJWT).To(ping))

	wsCont.Add(baseWs)

	// register routes
	userRep := repo.NewEgwUserRepository(db)
	userSvc := usecases.NewEgwUserService(userRep)

	user.NewEgwUserHandler(userSvc, wsCont)

	// product routes
	productRep := repo.NewEgwProductRepository(db)
	productSvc := usecases.NewEgwProductService(productRep)

	product.NewEgwProductHandler(productSvc, wsCont)

	http.Handle("/", wsCont)

	return fullSrv
}

func (s *Server) ListenAndServe(env string, domain string) error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func ping(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte("PONG"))
}

func (s *Server) WriteServiceErrorJson(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
	s.RequestLogger.Errorf("Service error: ", err)
	resp.WriteHeader(500)
	resp.WriteAsJson("internal server error")
}

func (s *Server) RecoverHandler(i interface{}, w http.ResponseWriter) {
	s.RequestLogger.Error("Server panic error: ", i)
	w.Write([]byte("Internal server error"))
}
