package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
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

	// TODO: show/redirect to "/play" page
	http.Redirect(w, r, fmt.Sprintf("/play/%s", file), 301)
}

func (s *Server) PlayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file"]

	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// TODO: this is very unsafe. Fix It!
	matches, _ := filepath.Glob(fmt.Sprintf("downloads/%s.*", base)) // TODO: get "downloads" from viper
	s.logger.Info("glob", log.Any("matches", matches))

	// TODO: find a better way
	video := ""
	thumb := ""
	for _, f := range matches {
		ext := filepath.Ext(f)
		switch ext {
		case ".mp4":
			fallthrough
		case ".webm":
			video = filepath.Base(f)
			break
		case ".jpg":
			fallthrough
		case ".png":
			thumb = filepath.Base(f)
			break
		default:
			s.logger.Error("unknown file ext", log.Any("ext", ext))
		}
	}

	fullData := map[string]interface{}{
		"Video": video,
		"Thumb": thumb,
	}
	files := []string{
		"templates/play.html",
		"templates/layout.html",
	}
	tpl := template.Must(template.New("base").ParseFiles(files...))
	s.render(w, r, tpl, fullData)
}
