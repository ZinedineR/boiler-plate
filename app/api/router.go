package api

import (
	"fmt"
	"ms-batch/internal/base/handler"
)

func (h *HttpServe) setupRouter() {

	//MailTemplate
	h.GuestRoute("POST", "/email-template", h.templateHandler.CreateMailTemplate)
	h.GuestRoute("GET", "/email-template", h.templateHandler.FindMailTemplate)
	h.GuestRoute("GET", "/email-template/:id", h.templateHandler.FindOneMailTemplate)
	h.GuestRoute("PUT", "/email-template/:id", h.templateHandler.UpdateMailTemplate)
	h.GuestRoute("DELETE", "/email-template/:id", h.templateHandler.DeleteMailTemplate)
}

func (h *HttpServe) UserRoute(method, path string, f handler.HandlerFnInterface) {
	switch method {
	case "GET":
		h.router.GET(path, h.base.UserRunAction(f))
	case "POST":
		h.router.POST(path, h.base.UserRunAction(f))
	case "PUT":
		h.router.PUT(path, h.base.UserRunAction(f))
	case "DELETE":
		h.router.DELETE(path, h.base.UserRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}

func (h *HttpServe) GuestRoute(method, path string, f handler.HandlerFnInterface) {
	switch method {
	case "GET":
		h.router.GET(path, h.base.GuestRunAction(f))
	case "POST":
		h.router.POST(path, h.base.GuestRunAction(f))
	case "PUT":
		h.router.PUT(path, h.base.GuestRunAction(f))
	case "DELETE":
		h.router.DELETE(path, h.base.GuestRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}
