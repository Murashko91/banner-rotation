package internalhttp

import (
	"net/http"

	"github.com/otus-murashko/banners-rotation/internal/app"
)

type Handler struct {
	app app.Application
}

func (h Handler) bannerRotationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBannerRotation(w, r, h.app)
	case http.MethodPost:
		addBannerRotation(w, r, h.app)
	case http.MethodDelete:
		deleteBannerRotation(w, r, h.app)
	default:
		handleNotExpecterRequest(w)
	}
}

func (h Handler) slotHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addSlot(w, r, h.app)
	default:
		handleNotExpecterRequest(w)
	}
}

func (h Handler) bannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addBanner(w, r, h.app)
	default:
		handleNotExpecterRequest(w)
	}
}

func (h Handler) groupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addGroup(w, r, h.app)
	default:
		handleNotExpecterRequest(w)
	}
}

func (h Handler) statHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		updateClickStat(w, r, h.app)
	default:
		handleNotExpecterRequest(w)
	}
}
