package api

import (
	"fmt"

	"boiler-plate/internal/base/handler"
)

func (h *HttpServe) setupProfileRouter() {
	h.GuestRoute("GET", "/profile", h.ProfileHandler.Find)
	h.GuestRoute("POST", "/profile", h.ProfileHandler.Create)
	h.GuestRoute("PUT", "/profile/:id", h.ProfileHandler.Update)
	h.GuestRoute("GET", "/profile/:id", h.ProfileHandler.Detail)
	h.GuestRoute("DELETE", "/profile/:id", h.ProfileHandler.Delete)
}

func (h *HttpServe) UserRoute(method, path string, f handler.HandlerFnInterface) {
	userRoute := h.router.Group("/api/v2")
	switch method {
	case "GET":
		userRoute.GET(path, h.base.UserRunAction(f))
	case "POST":
		userRoute.POST(path, h.base.UserRunAction(f))
	case "PUT":
		userRoute.PUT(path, h.base.UserRunAction(f))
	case "DELETE":
		userRoute.DELETE(path, h.base.UserRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}

func (h *HttpServe) GuestRoute(method, path string, f handler.HandlerFnInterface) {
	guestRoute := h.router.Group("/api/v2")
	switch method {
	case "GET":
		guestRoute.GET(path, h.base.GuestRunAction(f))
	case "POST":
		guestRoute.POST(path, h.base.GuestRunAction(f))
	case "PUT":
		guestRoute.PUT(path, h.base.GuestRunAction(f))
	case "DELETE":
		guestRoute.DELETE(path, h.base.GuestRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}
