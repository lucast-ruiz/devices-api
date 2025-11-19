package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	httpSwagger "github.com/swaggo/http-swagger"
	docs "github.com/lucast-ruiz/devices-api/internal/docs"

	"github.com/lucast-ruiz/devices-api/internal/api"
	"github.com/lucast-ruiz/devices-api/internal/repo"
	"github.com/lucast-ruiz/devices-api/internal/service"
)

// @title Devices API
// @version 1.0
// @description Devices management API
// @host localhost:8080
// @BasePath /
func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	deviceRepo := repo.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)
	handler := api.NewHandler(deviceService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	//Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	//Swagger
	docs.SwaggerInfo.BasePath = "/"
	
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), 
	))

	//api
	r.Mount("/", handler.Routes())

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}