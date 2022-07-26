package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jaredwarren/rpi_music/log"
)

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fullData := map[string]interface{}{}
	files := []string{
		"templates/index.html",
		"templates/layout.html",
	}
	tpl := template.Must(template.New("base").ParseFiles(files...))
	s.render(w, r, tpl, fullData)
}

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.httpError(w, fmt.Errorf("DownloadHandler|ParseForm|%w", err), http.StatusBadRequest)
		return
	}
	s.logger.Info("DownloadHandler", log.Any("form", r.PostForm))

	url := r.PostForm.Get("url")
	if url == "" {
		s.httpError(w, fmt.Errorf("need url"), http.StatusBadRequest)
		return
	}

	file, video, err := s.downloader.DownloadVideo(url)
	if err != nil {
		s.httpError(w, fmt.Errorf("DownloadHandler|downloadVideo|%w", err), http.StatusInternalServerError)
		return
	}
	tmb, err := s.downloader.DownloadThumb(video)
	if err != nil {
		s.logger.Warn("DownloadHandler|downloadThumb", log.Error(err))
		// ignore err
	}

	s.logger.Info("video downloaded", log.Any("file", file), log.Any("thumb", tmb))

	http.Redirect(w, r, "/", 301)
}
