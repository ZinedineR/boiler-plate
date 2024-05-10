package api

import (
	"fmt"

	"boiler-plate/internal/base/handler"
)

func (h *HttpServe) setupSettingsRouter() {
	h.GuestRoute("GET", "/settings", h.settingsHandler.FindSettings)
	h.GuestRoute("GET", "/general", h.settingsHandler.FindGeneral)
	h.GuestRoute("GET", "/logo", h.settingsHandler.FindLogo)
	h.GuestRoute("GET", "/password", h.settingsHandler.FindPassword)
	h.GuestRoute("PUT", "/general", h.settingsHandler.UpdateGeneral)
	h.GuestRoute("PUT", "/logo", h.settingsHandler.UpdateLogo)
	h.GuestRoute("PUT", "/password", h.settingsHandler.UpdatePassword)

	h.GuestRoute("GET", "/captcha", h.captchaHandler.GenerateCaptcha)
	h.GuestRoute("POST", "/login", h.loginHandler.Login)
}

func (h *HttpServe) setupAccountRouter() {
	h.UserRoute("POST", "/user/req-update-email", h.accountHandler.ReqUpdateEmail)
	h.UserRoute("POST", "/user/req-update-phone", h.accountHandler.ReqUpdatePhone)
	h.UserRoute("PUT", "/user/update-email", h.accountHandler.UpdateEmail)
	h.UserRoute("PUT", "/user/update-phone", h.accountHandler.UpdatePhone)
	h.GuestRoute("GET", "/user/contact-information/:id", h.accountHandler.GetContactInformation)

	h.UserRoute("PUT", "/account/information", h.accountHandler.UpdateAccountInformation)
}

func (h *HttpServe) setupRegistrationRouter() {
	h.GuestRoute("POST", "/register", h.registrationHandler.Registration)
}

func (h *HttpServe) setupVerifyRouter() {
	h.GuestRoute("POST", "/verify", h.verifyHandler.Verify)
}

func (h *HttpServe) setupInvestorCategoryRouter() {
	invCategoryRoute := h.router.Group("/api/v2/investor/category")

	invCategoryRoute.POST("/", h.base.UserRunAction(h.invCategoryHandler.Create))
	invCategoryRoute.GET("/:id", h.base.UserRunAction(h.invCategoryHandler.Detail))
	invCategoryRoute.DELETE("/:id", h.base.UserRunAction(h.invCategoryHandler.Delete))
	invCategoryRoute.GET("/", FilterMiddle(), h.base.UserRunAction(h.invCategoryHandler.List))
	invCategoryRoute.GET("/select", h.base.UserRunAction(h.invCategoryHandler.Select))
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
