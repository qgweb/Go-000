//+build wireinject

package main

import (
	"Go-000/Week04/internal/conf"
	xhttp "Go-000/Week04/internal/server/http"
	"Go-000/Week04/internal/service"
	"net/http"

	"github.com/google/wire"
)

func InitializeServer(config conf.Config, web *http.Server) *xhttp.WebServer {
	wire.Build(initHttp, initService)
	return &xhttp.WebServer{}
}

func initService(config conf.Config) *service.Service {
	svr := service.New(config)
	return svr
}

func initHttp(svr *service.Service, web *http.Server) *xhttp.WebServer {
	return xhttp.New(svr, web)
}
