package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"time"

	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"snippetbox.macnigor.net/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *models.UserModel
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error("error loading .env", "error", err)
		os.Exit(1)
	}

	//мы передаем через командную строку при запуске программы параметры
	//для запуска
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "4000"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("DATABASE_URL environment variable is not set")
	}

	flag.Parse()

	//создаем структурированный логер
	//&slog.HandlerOptions с его помощью можем отдавать логи нужного нам уровня здесь Info и выше

	db, err := openDB(dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connect database", "dsn", dsn)
	//закрываем соеденение с базой перед выходом из main
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	//инициализировали sessionManager добавили свои настройки
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	//инициализируем tls Config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//инициализировали  http.Server
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routers(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server on ", "addr", addr)

	//запустили сервер через srv
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
