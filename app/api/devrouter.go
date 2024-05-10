package api

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/base/app"
	"boiler-plate/pkg/server"
	"net/http"
)

func (h *HttpServe) setupDevRouter(conf *appconf.Config) {
	h.router.GET("/api/v2/health-check", h.base.GuestRunAction(func(ctx *app.Context) *server.ResponseInterface {
		return &server.ResponseInterface{
			Status: http.StatusOK,
			Data: map[string]interface{}{
				"status":  "ok",
				"service": conf.AppEnvConfig.AppName,
			},
		}
	}))
}
