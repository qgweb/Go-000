package http

import (
	"Go-000/Week04/api"
	"Go-000/Week04/internal/service"
	"Go-000/Week04/internal/xerror"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type WebServer struct {
	svr    *service.Service
	engine *http.Server
	webMux *http.ServeMux
}

func New(svr *service.Service, engine *http.Server) *WebServer {
	ws := &WebServer{
		svr:    svr,
		engine: engine,
		webMux: http.NewServeMux(),
	}
	ws.route()
	ws.engine.Handler = ws.webMux
	return ws
}

func response(code int, msg string, data interface{}) []byte {
	rsp, _ := json.Marshal(map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return rsp
}

func (ws *WebServer) bind(request *http.Request, val interface{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(val))
	if !refVal.CanSet() {
		return errors.New("can't set")
	}
	for i := 0; i < refVal.NumField(); i++ {
		fieldName := refVal.Type().Field(i).Tag.Get("json")
		if fieldName == "-" {
			continue
		}
		fieldName = strings.TrimSuffix(fieldName, ",omitempty")
		fieldVal := request.PostFormValue(fieldName)
		refVal.Field(i).Set(reflect.ValueOf(fieldVal))
	}
	return nil
}

func (ws *WebServer) route() {
	ws.webMux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		var loginRequest = &api.LoginReq{}
		if err := ws.bind(request, loginRequest); err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		info, err := ws.svr.SimpleLogin(ctx, loginRequest.UserName, loginRequest.Passwd)
		if err != nil {
			if xerror.IsNotFound(err) {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			if xerror.IsUnauthorized(err) {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			if xerror.IsInternal(err) {
				log.Printf("%+v", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(response(http.StatusOK, "ok", map[string]interface{}{
			"uid": info.ID,
		}))
	})
}

func (ws *WebServer) Close(ctx context.Context) {
	ws.svr.Close(ctx)
}
