package server

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaredwarren/plex_ytdl/downloader"
	"github.com/jaredwarren/rpi_music/log"
	"github.com/spf13/viper"
)

// Config provides basic configuration
type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Logger       log.Logger
}

// HTMLServer represents the web service that serves up HTML
type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
	logger log.Logger
}

// Start launches the HTML Server
func StartHTTPServer(cfg *Config) *HTMLServer {
	// Setup Context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init server
	s := New(cfg.Logger)

	// Setup Handlers
	r := mux.NewRouter()
	r.Use(s.loggingMiddleware)

	r.HandleFunc("/", s.HomeHandler).Methods("GET")
	r.HandleFunc("/download", s.DownloadHandler).Methods("POST")

	r.HandleFunc("/play/{file}", s.PlayHandler).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.PathPrefix("/video/").Handler(http.StripPrefix("/video/", http.FileServer(http.Dir(viper.GetString("player.song_root")))))
	r.PathPrefix("/thumb/").Handler(http.StripPrefix("/thumb/", http.FileServer(http.Dir(viper.GetString("player.thumb_root")))))

	// Create the HTML Server
	htmlServer := HTMLServer{
		logger: cfg.Logger,
		server: &http.Server{
			Addr:           cfg.Host,
			Handler:        r,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}

	// Start the listener
	htmlServer.wg.Add(1)
	go func() {
		cfg.Logger.Info("Starting HTTP server", log.Any("host", cfg.Host), log.Any("https", viper.GetBool("https")))
		if viper.GetBool("https") {
			htmlServer.server.ListenAndServeTLS("localhost.crt", "localhost.key")
		} else {
			htmlServer.server.ListenAndServe()
		}
		htmlServer.wg.Done()
	}()

	return &htmlServer
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Stop turns off the HTML Server
func (htmlServer *HTMLServer) StopHTTPServer() error {
	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	htmlServer.logger.Info("Stopping HTTP service...")

	// Attempt the graceful shutdown by closing the listener
	// and completing all inflight requests
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		// Looks like we timed out on the graceful shutdown. Force close.
		if err := htmlServer.server.Close(); err != nil {
			htmlServer.logger.Error("error stopping HTML service", log.Error(err))
			return err
		}
	}

	// Wait for the listener to report that it is closed.
	htmlServer.wg.Wait()
	htmlServer.logger.Info("HTTP service stopped")
	return nil
}

// Templates
type Server struct {
	logger     log.Logger
	downloader downloader.Downloader
}

func New(l log.Logger) *Server {
	return &Server{
		logger: l,
		downloader: &downloader.YoutubeDownloader{
			Logger: l,
		},
	}
}

// Render a template, or server error.
func (s *Server) render(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		s.logger.Error("template render error", log.Error(err))
		return
	}
	w.Write(buf.Bytes())
}

// Push the given resource to the client.
func (s *Server) push(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)
	if ok {
		err := pusher.Push(resource, nil)
		if err != nil {
			s.logger.Error("push error", log.Error(err))
		}
		return
	}
}
