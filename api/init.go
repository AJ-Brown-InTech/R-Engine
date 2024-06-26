package api

import (
	"Engine/storage"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)


func InitHandlers(router *chi.Mux, db *storage.DB) {
	
	// authenticated routes
	apiRouter := chi.NewRouter()
	apiRouter.Use(SessionMiddleware)
	apiRouter.Use(RequestMiddleware)

	router.Post("/register", RegisterAccount(db))
	router.Mount("/api/v1/", apiRouter)



	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting server, listening on port %s", port)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	
}