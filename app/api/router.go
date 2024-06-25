package api

import (
	"fmt"

	"boiler-plate/internal/base/handler"
)

func (h *HttpServe) setupUsersRouter() {
	h.GuestRoute("GET", "/users", h.UsersHandler.Find)
	h.GuestRoute("POST", "/users", h.UsersHandler.Create)
	h.GuestRoute("PUT", "/users/:id", h.UsersHandler.Update)
	h.GuestRoute("GET", "/users/:id", h.UsersHandler.Detail)
	h.GuestRoute("DELETE", "/users/:id", h.UsersHandler.Delete)
	h.GuestRoute("GET", "/users/:id/submissions", h.SubmissionsHandler.FindByUser)

	h.GuestRoute("GET", "/submissions", h.SubmissionsHandler.Find)
	h.GuestRoute("POST", "/submissions", h.SubmissionsHandler.Create)
	h.GuestRoute("GET", "/submissions/:id", h.SubmissionsHandler.Detail)
	h.GuestRoute("DELETE", "/submissions/:id", h.SubmissionsHandler.Delete)
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
		guestRoute.GET(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "POST":
		guestRoute.POST(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "PUT":
		guestRoute.PUT(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "DELETE":
		guestRoute.DELETE(path, AuthMiddle(), h.base.GuestRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}
