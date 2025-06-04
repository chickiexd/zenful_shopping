package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/chickiexd/zenful_shopping/internal/handler"
	// "github.com/chickiexd/zenful_shopping/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config  config
	handler handler.Handler
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
	// maxOpenConns int
	// maxIdleConns int
	// maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/ingredients", func(r chi.Router) {
			r.Get("/", app.handler.Ingredients.GetAll)
			r.Post("/create", app.handler.Ingredients.Create)
			r.Post("/add", app.handler.Ingredients.AddToShoppingList)
		})
		r.Route("/measurement_units", func(r chi.Router) {
			r.Get("/", app.handler.MeasurementUnits.GetAll)
			r.Post("/", app.handler.MeasurementUnits.Create)
		})
		r.Route("/food_groups", func(r chi.Router) {
			r.Get("/", app.handler.FoodGroups.GetAll)
			r.Post("/", app.handler.FoodGroups.Create)
		})
		r.Route("/recipes", func(r chi.Router) {
			r.Get("/", app.handler.Recipes.GetAll)
			r.Post("/", app.handler.Recipes.Create)
			r.Post("/parse", app.handler.ChatGPT.ParseRecipe)
			r.Post("/add/{id}", app.handler.Recipes.AddToShoppingList)
			r.Post("/remove", app.handler.Recipes.RemoveFromShoppingList)
		})
		r.Route("/images", func(r chi.Router) {
			r.Get("/{filename}", app.handler.Images.Get)
		})
		r.Route("/shopping_lists", func(r chi.Router) {
			r.Get("/", app.handler.ShoppingList.GetAll)
			r.Post("/remove_item", app.handler.ShoppingList.RemoveItemFromShoppingList)
			r.Post("/remove_all_items", app.handler.ShoppingList.RemoveAllItemsFromShoppingList)
		})
		r.Route("/sync", func(r chi.Router) {
			r.Get("/", app.handler.KeepSync.SyncShoppingLists)
		})
	})

	//recipes
	//users
	//ingredients
	//auth

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
