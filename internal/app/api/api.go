package api

import (
	"booking/configs"
	"net/http"

	memberhandler "booking/internal/app/api/handler/member"
	memberRepository "booking/internal/app/repositories/member"
	memberServices "booking/internal/app/services/member"

	"booking/internal/app/db"

	"booking/internal/pkg/glog"
	"booking/internal/pkg/health"
	"booking/internal/pkg/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	// InfraConns holds infrastructure services connections like MongoDB, Redis, Kafka,...
	InfraConns struct {
		Databases db.Connections
	}

	middlewareFunc = func(http.HandlerFunc, *configs.ErrorMessage) http.HandlerFunc
	route          struct {
		path        string
		method      string
		handler     http.HandlerFunc
		middlewares []middlewareFunc
	}
)

const (
	get    = http.MethodGet
	post   = http.MethodPost
	put    = http.MethodPut
	delete = http.MethodDelete
)

// Init init all handlers
func Init(conns *configs.Configs, em configs.ErrorMessage) (http.Handler, error) {
	logger := glog.New()

	// declare variable to pointer service class
	var memberRepo memberServices.Repository

	switch conns.Database.Type {
	case db.TypeMongoDB:
		s, err := configs.Dial(&conns.Database.Mongo, logger)
		if err != nil {
			logger.Panicf("failed to dial to target server, err: %v", err)
		}
		memberRepo = memberRepository.NewMongoRepository(s)

	default:
		panic("database type not supported: " + conns.Database.Type)
	}

	memberLogger := logger.WithField("package", "member")
	memberSrv := memberServices.NewService(conns, &em, memberRepo, memberLogger)
	memberHandler := memberhandler.New(conns, &em, memberSrv, memberLogger)

	routes := []route{
		// infra
		route{
			path:    "/readiness",
			method:  get,
			handler: health.Readiness().ServeHTTP,
		},
		// services
		// member
		route{
			path:        "/api/v1/member/{id:[a-z0-9-\\-]+}",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     memberHandler.Get,
		},
		route{
			path:        "/api/v1/member",
			method:      post,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     memberHandler.InsertMember,
		},
		route{
			path:        "/api/v1/member",
			method:      put,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     memberHandler.UpdateMemberByID,
		},
		route{
			path:    "/login",
			method:  post,
			handler: memberHandler.Login,
		},
	}

	loggingMW := middleware.Logging(logger.WithField("package", "middleware"))
	r := mux.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.StatusResponseWriter)
	r.Use(loggingMW)
	r.Use(handlers.CompressHandler)

	for _, rt := range routes {
		h := rt.handler
		for _, mdw := range rt.middlewares {
			h = mdw(h, &em)
		}
		r.Path(rt.path).Methods(rt.method).HandlerFunc(h)
	}

	return r, nil
}

// Close close all underlying connections
func (c *InfraConns) Close() {
	c.Databases.Close()
}
