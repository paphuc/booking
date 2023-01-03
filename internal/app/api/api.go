package api

import (
	"net/http"
	"booking/configs"
	memberhandler "booking/internal/app/api/handler/member"
	"booking/internal/app/db"
	"booking/internal/app/member"
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

	middlewareFunc = func(http.HandlerFunc) http.HandlerFunc
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
	var memberRepo member.Repository

	switch conns.Database.Type {
	case db.TypeMongoDB:
		s, err := configs.Dial(&conns.Database.Mongo, logger)
		if err != nil {
			logger.Panicf("failed to dial to target server, err: %v", err)
		}
		memberRepo = member.NewMongoRepository(s)
		

	default:
		panic("database type not supported: " + conns.Database.Type)
	}

	memberLogger := logger.WithField("package", "member")
	memberSrv := member.NewService(memberRepo, memberLogger)
	memberHandler := memberhandler.New(memberSrv, memberLogger)

	routes := []route{
		// infra
		route{
			path:    "/readiness",
			method:  get,
			handler: health.Readiness().ServeHTTP,
		},
		// services
		route{
			path:    "/api/v1/member/{id:[a-z0-9-\\-]+}",
			method:  get,
			handler: memberHandler.Get,
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
			h = mdw(h)
		}
		r.Path(rt.path).Methods(rt.method).HandlerFunc(h)
	}

	return r, nil
}

// Close close all underlying connections
func (c *InfraConns) Close() {
	c.Databases.Close()
}
