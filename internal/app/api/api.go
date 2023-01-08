package api

import (
	"booking/configs"
	"net/http"

	memberhandler "booking/internal/app/api/handler/member"
	memberServices "booking/internal/app/services/member"
	memberRepository "booking/internal/app/repositories/member"

	tablehandler "booking/internal/app/api/handler/table"
	tableServices "booking/internal/app/services/table"
	tableRepository "booking/internal/app/repositories/table"

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

	// declare variable to pointer repository class
	var memberRepo memberServices.Repository
	var tableRepo tableServices.Repository

	switch conns.Database.Type {
	case db.TypeMongoDB:
		s, err := configs.Dial(&conns.Database.Mongo, logger)
		if err != nil {
			logger.Panicf("failed to dial to target server, err: %v", err)
		}
		memberRepo = memberRepository.NewMongoRepository(s)
		tableRepo = tableRepository.NewMongoRepository(s)

	default:
		panic("database type not supported: " + conns.Database.Type)
	}

	memberLogger := logger.WithField("package", "member")
	memberSrv := memberServices.NewService(conns, &em, memberRepo, memberLogger)
	memberHandler := memberhandler.New(conns, &em, memberSrv, memberLogger)

	tableLogger := logger.WithField("package", "table")
	tableSrv := tableServices.NewService(conns, &em, tableRepo, tableLogger)
	tableHandler := tablehandler.New(conns, &em, tableSrv, tableLogger)

	routes := []route{
		// infra
		route{
			path:    "/readiness",
			method:  get,
			handler: health.Readiness().ServeHTTP,
		},
		// services
		// api member
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
		// api table
		route{
			path:        "/api/v1/table",
			method:      post,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     tableHandler.InsertTable,
		},
		route{
			path:        "/api/v1/table",
			method:      put,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     tableHandler.UpdateTableByID,
		},
		route{
			path:        "/api/v1/table-delete",
			method:      put,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     tableHandler.DeleteTable,
		},
		// api login
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
