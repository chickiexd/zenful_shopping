package main

import (
	"github.com/chickiexd/zenful_shopping/internal/db"
	"github.com/chickiexd/zenful_shopping/internal/env"
	"github.com/chickiexd/zenful_shopping/internal/handler"
	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/chickiexd/zenful_shopping/internal/service"
	"github.com/chickiexd/zenful_shopping/internal/store"
)

const version = "1.1.3"

//	@title			Zenful Shopping API
//	@description	This is the API for Zenful Shopping, a platform to manage your recipes, ingredients and shopping lists.

//	@contact.name	chickie
//	@contact.url	chickiexd.com
//	@contact.email	contact@chickiexd.com

//	@license.name	MIT License
//	@license.url	https://opensource.org/licenses/MIT

//	@BasePath	/v1

func main() {

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("API_URL", "http://localhost:8080"),
		db: dbConfig{
			host:     env.GetString("DB_HOST", "localhost"),
			user:     env.GetString("DB_USER", "user"),
			password: env.GetString("DB_PASSWORD", "adminpassword"),
			dbname:   env.GetString("DB_NAME", "zenful_shopping"),
			port:     env.GetString("DB_PORT", "9432"),
			// maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			// maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			// maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
		env: env.GetString("ENV", "dev"),
	}

	logger.Init()
	defer logger.Sync()

	db, err := db.New(cfg.db.host, cfg.db.user, cfg.db.password, cfg.db.dbname, cfg.db.port)
	if err != nil {
		logger.Log.Panic(err)
	}
	logger.Log.Info("db connection established")

	store := store.NewStorage(db)
	service := service.NewService(&store)
	handler := handler.NewHandler(&service)

	app := &application{
		config:  cfg,
		handler: handler,
	}

	mux := app.mount()

	logger.Log.Fatal(app.run(mux))
}
