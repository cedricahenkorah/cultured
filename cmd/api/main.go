package main

import (
	"context"
	"cultured/internal/pg"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	reviews        *pg.ReviewModel
	users          *pg.UserModel
	sessionManager *scs.SessionManager
	config         config
}

func main() {
	var cfg config

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	flag.IntVar(&cfg.port, "port", 4000, "HTTP network port")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "Environment(development|staging|production)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if cfg.db.dsn == "" {
		errorLog.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := openDB(cfg.db.dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		reviews:        &pg.ReviewModel{DB: db},
		users:          &pg.UserModel{DB: db},
		sessionManager: sessionManager,
		config:         cfg,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("cultured %s is starting on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
