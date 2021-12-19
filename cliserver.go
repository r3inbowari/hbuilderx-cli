package meiwobuxing

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/mux"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"github.com/wuwenbao/gcors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var BiliServer *Server

type Server struct {
	router *mux.Router
	s      *http.Server
}

func CLIApplication() {
	Log.Info("[BCS] MEIWOBUXING CLI PACKAGER is running")
	BiliServer = NewServer()
	BiliServer.Map("/version", HandleVersion)
	BiliServer.Map("/file", FileUpload, http.MethodPost)
	BiliServer.Map("/pack", HandlePackRequest, http.MethodPost)
	BiliServer.Map("/state/{traceID}", HandlePackState, http.MethodGet)
	BiliServer.Map("/manifest", HandleManifest, http.MethodPost)
	BiliServer.Map("/file/{id}", FileDownload, http.MethodGet)
	BiliServer.Map("/export", ExportElements, http.MethodPost)
	BiliServer.Map("/callback", CallbackTest, http.MethodGet)
	BiliServer.Map("/login", HandleLogin, http.MethodPost)
	BiliServer.router.Use(loggingMiddleware)
	BiliServer.router.Use(authMiddleware())
	AppInfo()
	err := BiliServer.start()
	if strings.HasSuffix(err.Error(), "normally permitted.") || strings.Index(err.Error(), "bind") != -1 {
		Log.WithFields(logrus.Fields{"err": err.Error()}).Error("[BCS] only one usage of each socket address is normally permitted.")
		Log.Info("[BCS] EXIT 1002")
		time.Sleep(time.Second * 5)
		os.Exit(1002)
	}
	// goroutine block here not need sleep
	Log.WithFields(logrus.Fields{"err": err.Error()}).Info("[BCS] service will be terminated soon")
	time.Sleep(time.Second * 10)
}

//func authMiddleware1(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		Log.Info("[BSC] route " + r.RemoteAddr + " " + r.RequestURI)
//		next.ServeHTTP(w, r)
//	})
//}

func Shutdown(ctx context.Context) {
	BiliServer.Shutdown(ctx)
}

func NewServer() *Server {
	r := mux.NewRouter()
	Log.Info("[BSC] global CORS open")

	cors := gcors.New(
		r,
		gcors.WithOrigin("*"),
		gcors.WithMethods("POST, GET, PUT, DELETE, OPTIONS"),
		gcors.WithHeaders("Authorization"),
	)

	retServer := &http.Server{
		Addr:    GetConfig(false).APIAddr,
		Handler: cors,
	}
	return &Server{router: r, s: retServer}
}

func (s *Server) start() error {
	Log.Info("[BCS] listened on " + GetConfig(false).APIAddr)
	if GetConfig(false).CaCert != "" && GetConfig(false).CaKey != "" {
		certPath := Up.RunPath + "/" + GetConfig(false).CaCert
		keyPath := Up.RunPath + "/" + GetConfig(false).CaKey
		_, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			Log.Error("[BSC] please check your cert path whether is right")
			return err
		}
		Log.Info("[BSC] tls enabled")
		return s.s.ListenAndServeTLS(certPath, keyPath)
	} else {
		return s.s.ListenAndServe()
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.s != nil {
		Log.Info("[BSC] releasing server now...")
		err := s.s.Shutdown(ctx)
		if err != nil {
			Log.Error("[BSC] shutdown failed")
			Log.Info("[BCS] EXIT 1002")
			time.Sleep(time.Second * 5)
			os.Exit(1011)
		}
		Log.Info("[BSC] release completed")
	}
}

func (s *Server) Map(path string, f func(http.ResponseWriter,
	*http.Request), method ...string) *Server {
	if len(method) == 1 {
		Log.Info("[BSC] add route path [" + method[0] + "] -> " + path)
		s.router.HandleFunc(path, f).Methods(method[0])
	} else {
		Log.Info("[BSC] add route path [ALL] -> " + path)
		s.router.HandleFunc(path, f)
	}
	return s
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Info("[BSC] route " + r.RemoteAddr + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware() mux.MiddlewareFunc {
	if GetConfig(false).JwtEnable {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI == "/login" || r.RequestURI == "/version" {

				} else {
					auth := r.Header.Get("Authorization")
					sa := strings.Split(auth, " ")
					if len(sa) != 2 {
						ResponseCommon(w, "401", "request failed", 1, http.StatusOK, 6401)
						return
					}
					err := CheckToken(sa[1])
					if err != nil {
						ResponseCommon(w, err.Error(), "request failed", 1, http.StatusOK, 6401)
						return
					}
				}
				next.ServeHTTP(w, r)
			})
		}
	} else {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}
}

func AppInfo() {
	var RunningMode = "server"
	if Up == nil || len(Up.ReleaseTag) != 40 {
		Log.Error("[MAIN] illegally generated program")
		time.Sleep(time.Second * 3)
		os.Exit(33)
	}
	Log.Blue(" ________  ________  ___  ________  ________  ___     ")
	Log.Blue("|\\   ____\\|\\   __  \\|\\  \\|\\   ____\\|\\   __  \\|\\  \\          PACKAGER #UNOFFICIAL " + Up.ReleaseTag[:7] + "..." + Up.ReleaseTag[33:])
	Log.Blue("\\ \\  \\___|\\ \\  \\|\\  \\ \\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\         -... .. .-.. .. -.-. --- .. -. " + Up.VersionStr)
	Log.Blue(" \\ \\  \\    \\ \\   __  \\ \\  \\ \\  \\    \\ \\   __  \\ \\  \\        Running: CLI Server" + " by cyt(r3inbowari)")
	if RunningMode == "server" {
		Log.Blue("  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\       Listened: " + GetConfig(false).APIAddr)
	} else {
		Log.Blue("  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\       Port: UNSUPPORTED")
	}
	Log.Blue("   \\ \\_______\\ \\__\\ \\__\\ \\__\\ \\_______\\ \\__\\ \\__\\ \\__\\      PID: " + strconv.Itoa(os.Getpid()))
	Log.Blue("    \\|_______|\\|__|\\|__|\\|__|\\|_______|\\|__|\\|__|\\|__|      Built on " + Up.BuildTime)
	Log.Blue("")
}
