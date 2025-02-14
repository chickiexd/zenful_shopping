package main

import (
	"log"
	"zenful_shopping_backend/internal/db"
	"zenful_shopping_backend/internal/env"
	"zenful_shopping_backend/internal/store"
	"zenful_shopping_backend/internal/service"
	"zenful_shopping_backend/internal/handler"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
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
	}

	db, err := db.New(cfg.db.host, cfg.db.user, cfg.db.password, cfg.db.dbname, cfg.db.port)
	if err != nil {
		log.Panic(err)
	}
	log.Println("db connection established")

	service := service.NewService()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
