package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SibHelly/url_shortnener/api"
	"github.com/SibHelly/url_shortnener/api/middleware"
	"github.com/SibHelly/url_shortnener/internals/app/db"
	"github.com/SibHelly/url_shortnener/internals/app/handlers"
	urlprocessor "github.com/SibHelly/url_shortnener/internals/app/processors"
	"github.com/SibHelly/url_shortnener/internals/cfg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server
}

func (server *Server) Serve() {
	log.Println("Starting server")
	// соединие с бд
	var err error
	server.db, err = pgxpool.New(server.ctx, server.config.GetDbString())
	if err != nil {
		log.Fatalln(err)
	}

	urlStorage := db.NewUrlStorage(server.db)
	urlprocessor := urlprocessor.NewUrlProcessor(urlStorage)
	urlHandler := handlers.NewUrlHandler(urlprocessor)

	routes := api.CreateRoutes(urlHandler)
	routes.Use(middleware.RequestLog)
	//  роуты
	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	log.Println("Server started")

	err = server.srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	return
}

func (server *Server) Shutdown() {
	log.Println("Server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()
	defer func() {
		cancel()
	}()
	var err error
	if err = server.srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
