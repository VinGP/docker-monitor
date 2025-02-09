package handlers

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/api/response"
	"backend/pkg/logger/sl"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"

	docs "backend/docs"

	"log/slog"
	"net/http"

	mwLogger "backend/internal/http/middleware/logger"
)

const corsMaxAge = 300 // Maximum value not ignored by any of major browsers

// NewRouter
// @description     API for Docker Monitor.
// @version         1.0
// @title           Docker Monitor API
// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT
func NewRouter(r chi.Router, log *slog.Logger, csService *service.ContainerStatusService) {
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Real-IP"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           corsMaxAge,
	}))
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/api/docs/*", func(w http.ResponseWriter, r *http.Request) {
		baseURL := r.Host
		docs.SwaggerInfo.Host = baseURL
		httpSwagger.Handler(
			httpSwagger.URL("/api/docs/doc.json"), // The URL pointing to API definition
		).ServeHTTP(w, r)
	})

	csr := NewContainerStatusResource(csService)
	r.Route("/api/container_status", func(r chi.Router) {
		r.Get("/", csr.GetAllContainerStatuses)
		r.Delete("/", csr.DeleteAllContainerStatuses)
		r.Post("/", csr.CreateContainerStatus)
	})
}

type ContainerStatusResource struct {
	cs *service.ContainerStatusService
}

func NewContainerStatusResource(cs *service.ContainerStatusService) *ContainerStatusResource {
	return &ContainerStatusResource{cs}
}

// GetAllContainerStatuses
// @Tags container_status
// @Summary Get all container status
// @Description Get all container status
// @Produce json
// @Success 200 {array} model.ContainerStatus
// @Failure 500 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/container_status [get].
func (csr *ContainerStatusResource) GetAllContainerStatuses(w http.ResponseWriter, r *http.Request) {
	containers, err := csr.cs.GetAll()
	if err != nil {
		slog.Error("GetAllContainerStatuses", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		renderJson(w, response.Error(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	renderJson(w, &containers)
}

// DeleteAllContainerStatuses
// @Tags container_status
// @Summary Delete all container status
// @Description Delete all container status
// @Produce json
// @Router /api/container_status [delete].
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure 400 {object} response.Response
func (csr *ContainerStatusResource) DeleteAllContainerStatuses(w http.ResponseWriter, r *http.Request) {
	err := csr.cs.DeleteAll()
	if err != nil {
		slog.Error("DeleteAllContainerStatuses", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		renderJson(w, response.Error(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	renderJson(w, response.OK())
}

// CreateContainerStatus
// @Tags container_status
// @Summary Create container status
// @Description Create container status
// @Param container_status body model.ContainerStatus true "Container status".
// @Router /api/container_status [post].
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
func (csr *ContainerStatusResource) CreateContainerStatus(w http.ResponseWriter, r *http.Request) {
	var req model.ContainerStatus

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		renderJson(w, response.Error(err.Error()))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		renderJson(w, response.ValidationError(err.(validator.ValidationErrors))) //nolint:errcheck,errorlint // ignore
		return
	}

	err := csr.cs.UpsertContainerStatus(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		renderJson(w, response.Error(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	renderJson(w, response.OK())
}

func renderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("renderJson", sl.Err(err))
	}
}
