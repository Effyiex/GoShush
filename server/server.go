package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type HostConfig struct {
	Port        string  `toml:"port"`
	Label       string  `toml:"label"`
	ShortLabel  string  `toml:"short_label"`
	FrontColor  [3]byte `toml:"front_color"`
	AccentColor [3]byte `toml:"accent_color"`
	BackColor   [3]byte `toml:"back_color"`
	DataFolder  string  `toml:"data_folder"`
}

type WebAppIcon struct {
	Path     string `json:"src"`
	MimeType string `json:"type"`
	Scale    string `json:"sizes"`
}

type WebAppManifest struct {
	Name            string       `json:"name"`
	StartUrl        string       `json:"start_url"`
	BackgroundColor string       `json:"background_color"`
	ThemeColor      string       `json:"theme_color"`
	Display         string       `json:"display"`
	CaptureLinks    string       `json:"capture_links"`
	Icons           []WebAppIcon `json:"icons"`
}

func NewAppManifest(cfg HostConfig) WebAppManifest {
	return WebAppManifest{
		Name:            cfg.Label,
		StartUrl:        "/offline?installed",
		BackgroundColor: "#" + fmt.Sprintf("%02x%02x%02x", cfg.BackColor[0], cfg.BackColor[1], cfg.BackColor[2]),
		ThemeColor:      "#" + fmt.Sprintf("%02x%02x%02x", cfg.FrontColor[0], cfg.FrontColor[1], cfg.FrontColor[2]),
		Display:         "standalone",
		CaptureLinks:    "existing-client-navigate",
		Icons: []WebAppIcon{{
			Path:     "favicon.png",
			Scale:    "256x256",
			MimeType: "image/png",
		}, {
			Path:     "favicon.ico",
			Scale:    "256x256",
			MimeType: "image/x-icon",
		}},
	}
}

type InterceptFunc = func(data []byte) []byte

type Server struct {
	*mux.Router
	Users                UserManager
	Chats                ChatManager
	Config               HostConfig
	AppManifest          WebAppManifest
	Intercepts           map[string]InterceptFunc
	GeneralHTMLIntercept InterceptFunc
	FileQueryRoutes      []string
}

func NewServer(cfg HostConfig) *Server {
	if _, err := os.Stat(cfg.DataFolder); os.IsNotExist(err) {
		os.Mkdir(cfg.DataFolder, os.ModePerm)
	}
	users := BootUserManager(cfg.DataFolder)
	chats := BootChatManager(cfg.DataFolder)
	srv := &Server{
		Router:          mux.NewRouter(),
		Users:           users,
		Chats:           chats,
		Config:          cfg,
		AppManifest:     NewAppManifest(cfg),
		Intercepts:      map[string]InterceptFunc{},
		FileQueryRoutes: make([]string, 0),
	}
	srv.GeneralHTMLIntercept = srv.NewGeneralHTMLIntercept()
	srv.RegisterBackRoutes()
	srv.RegisterFrontRoutes()
	return srv
}

func (srv *Server) Run() {
	fmt.Println("Launching CommandLine-Listener.")
	go srv.CommandLineListener()
	fmt.Println("Hosting on \"127.0.0.1:" + srv.Config.Port + "\".")
	http.ListenAndServe(":"+srv.Config.Port, srv)
}

func (srv *Server) HandleIntercept(route string, interceptor InterceptFunc) {
	srv.Intercepts[route] = interceptor
}

func (srv *Server) Dispose(status int) {
	fmt.Println("Saving Users...")
	srv.Users.Save(srv.Config.DataFolder)
	fmt.Println("Saving Chats...")
	srv.Chats.Save(srv.Config.DataFolder)
	fmt.Println("Shutting down!")
	os.Exit(status)
}

var mimeTypeMapping = map[string]string{
	".css":  "text/css",
	".js":   "text/javascript",
	".png":  "image/png",
	".ico":  "image/x-icon",
	".html": "text/html",
}

func (srv *Server) PeekMimeType(path string) string {
	mimeType, specific := mimeTypeMapping[filepath.Ext(path)]
	if specific {
		return mimeType
	} else {
		return "application/octet-stream"
	}
}

const WEB_FILES_ROOT = "web/"

func (srv *Server) NewGeneralHTMLIntercept() InterceptFunc {
	return func(data []byte) []byte {
		html := string(data)
		html = strings.ReplaceAll(html, "%appname%", srv.Config.Label)
		return []byte(html)
	}
}

func (srv *Server) HandleFileQuery(route string, path string) {

	srv.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile(path)
		if interceptor, ok := srv.Intercepts[route]; ok {
			data = interceptor(data)
		}
		if strings.HasSuffix(path, ".html") {
			data = srv.GeneralHTMLIntercept(data)
		}
		w.Header().Add("Content-Type", srv.PeekMimeType(path))
		w.Write(data)
	}).Methods("GET")

	fmt.Println("Registered route: \"" + route + "\", for path: \"" + path + "\".")

}

func (srv *Server) RegisterBackRoutes() {
	srv.Users.RegisterOn(*srv)
	srv.Chats.RegisterOn(*srv)
}

func (srv *Server) WalkWebFilesRoot() {

	filepath.Walk(WEB_FILES_ROOT, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		path = strings.ReplaceAll(path, "\\", "/")

		var route = path
		route, _ = strings.CutPrefix(route, WEB_FILES_ROOT)
		route, _ = strings.CutSuffix(route, ".html")
		route = "/" + route

		var alr_exists = false
		for _, _route := range srv.FileQueryRoutes {
			if _route == route {
				alr_exists = true
				break
			}
		}
		if !alr_exists {
			srv.FileQueryRoutes = append(srv.FileQueryRoutes, route)
			srv.HandleFileQuery(route, path)
		}

		return nil

	})

}

func (srv *Server) RegisterFrontRoutes() {

	srv.WalkWebFilesRoot()
	srv.HandleFileQuery("/", "web/login.html")

	srv.HandleIntercept("/css/master.css", func(data []byte) []byte {
		script := string(data)
		script = strings.ReplaceAll(
			script,
			"--backcolor_r: none",
			"--backcolor_r: "+fmt.Sprint(int(srv.Config.BackColor[0])),
		)
		script = strings.ReplaceAll(
			script,
			"--backcolor_g: none",
			"--backcolor_g: "+fmt.Sprint(int(srv.Config.BackColor[1])),
		)
		script = strings.ReplaceAll(
			script,
			"--backcolor_b: none",
			"--backcolor_b: "+fmt.Sprint(int(srv.Config.BackColor[2])),
		)
		script = strings.ReplaceAll(
			script,
			"--accentcolor_r: none",
			"--accentcolor_r: "+fmt.Sprint(int(srv.Config.AccentColor[0])),
		)
		script = strings.ReplaceAll(
			script,
			"--accentcolor_g: none",
			"--accentcolor_g: "+fmt.Sprint(int(srv.Config.AccentColor[1])),
		)
		script = strings.ReplaceAll(
			script,
			"--accentcolor_b: none",
			"--accentcolor_b: "+fmt.Sprint(int(srv.Config.AccentColor[2])),
		)
		script = strings.ReplaceAll(
			script,
			"--frontcolor_r: none",
			"--frontcolor_r: "+fmt.Sprint(int(srv.Config.FrontColor[0])),
		)
		script = strings.ReplaceAll(
			script,
			"--frontcolor_g: none",
			"--frontcolor_g: "+fmt.Sprint(int(srv.Config.FrontColor[1])),
		)
		script = strings.ReplaceAll(
			script,
			"--frontcolor_b: none",
			"--frontcolor_b: "+fmt.Sprint(int(srv.Config.FrontColor[2])),
		)
		return []byte(script)
	})

	srv.HandleFunc("/app.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/manifest+json")
		data, _ := json.Marshal(srv.AppManifest)
		w.Write(data)
	}).Methods("GET")

}
