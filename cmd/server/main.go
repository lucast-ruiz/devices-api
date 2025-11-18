package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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
	docs.SwaggerInfo.BasePath = "/"
	
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), 
	))

	r.Mount("/", handler.Routes())

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}