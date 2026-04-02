package api

import (
	"errors"
	"net/http"

	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/nirrax/url_shortener/internal/service"
	"github.com/nirrax/url_shortener/internal/utils"
)

func (app *App) NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /url/{url}", app.ShowUrlHandler)
	router.HandleFunc("GET /{url}", app.RedirectUrlHandler)
	router.HandleFunc("POST /", app.CreateUrlHandler)
	router.HandleFunc("DELETE /{url}", app.DeleteUrlHandler)

	return router
}

func (app *App) ShowUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := utils.ExtractUrl(r)
	if err != nil {
		utils.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	urlDTO, err := app.UrlService().GetUrl(url)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUrlNotFound):
			utils.SendJsonError(w, http.StatusNotFound, err)
		default:
			utils.SendJsonError(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.SendJsonResponse(w, http.StatusOK, urlDTO)
}

func (app *App) RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := utils.ExtractUrl(r)
	if err != nil {
		utils.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	urlDTO, err := app.UrlService().GetUrl(url)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUrlNotFound):
			utils.SendJsonError(w, http.StatusNotFound, err)
		default:
			utils.SendJsonError(w, http.StatusInternalServerError, err)
		}

		return
	}

	http.Redirect(w, r, urlDTO.LongUrl, http.StatusMovedPermanently)
}

func (app *App) CreateUrlHandler(w http.ResponseWriter, r *http.Request) {
	payload := entity.UrlRequest{}

	err := utils.DecodeJSONBody(r, &payload)
	if err != nil {
		utils.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	urlDto, err := app.UrlService().CreateUrl(payload.Url)
	if err != nil {
		utils.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendJsonResponse(w, http.StatusCreated, urlDto)
}

func (app *App) DeleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := utils.ExtractUrl(r)
	if err != nil {
		utils.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	err = app.UrlService().DeleteUrl(url)
	if err != nil {
		utils.SendJsonError(w, http.StatusNotFound, err)
		return
	}

	utils.SendJsonResponse(w, http.StatusOK, struct{}{})
}
