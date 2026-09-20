package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *pgxpool.Pool
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if *dsn == "" {
		errorLog.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("cultured is starting on %s", *addr)
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
