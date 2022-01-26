package meiwobuxing

import (
	"context"
	. "github.com/r3inbowari/zlog"
	"github.com/r3inbowari/zserver"
	"net/http"
	"os"
	"strconv"
	"time"
)

var s *zserver.Server

func CLIApplication() {
	Log.Info("[BCS] MEIWOBUXING CLI PACKAGER is running")
	s = zserver.DefaultServer(zserver.Options{
		Log:  &Log.Logger,
		Addr: GetConfig(false).APIAddr,
	})

	s.Map("/version", HandleVersion)
	s.Map("/file", FileUpload, http.MethodPost)
	s.Map("/pack", HandlePackRequest, http.MethodPost)
	s.Map("/state/{traceID}", HandlePackState, http.MethodGet)
	s.Map("/manifest", HandleManifest, http.MethodPost)
	s.Map("/file/{id}", FileDownload, http.MethodGet)
	s.Map("/export", ExportElements, http.MethodPost)
	s.Map("/callback", CallbackTest, http.MethodGet)
	s.Map("/login", HandleLogin, http.MethodPost)

	s.Map("/path", HandleAddPathMapping, http.MethodPost)
	s.Map("/path", HandleDelPathMapping, http.MethodDelete)
	s.Map("/paths", HandleAllPathMapping, http.MethodGet)

	s.Map("/project", HandleUploadProject, http.MethodPost)

	s.R.Use(s.AuthMiddlewareBuilder(func(token string, r *http.Request) error {
		return CheckToken(token)
	}))
	AppInfo()
	s.Start()
}

func AppInfo() {
	if Up == nil || len(Up.ReleaseTag) != 40 {
		Log.Error("[MAIN] illegally generated program")
		time.Sleep(time.Second * 3)
		os.Exit(33)
	}
	Log.Blue(" ________  ________  ___  ________  ________  ___     ")
	Log.Blue("|\\   ____\\|\\   __  \\|\\  \\|\\   ____\\|\\   __  \\|\\  \\          PACKAGER #UNOFFICIAL " + Up.ReleaseTag[:7] + "..." + Up.ReleaseTag[33:])
	Log.Blue("\\ \\  \\___|\\ \\  \\|\\  \\ \\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\         -... .. .-.. .. -.-. --- .. -. " + Up.VersionStr)
	Log.Blue(" \\ \\  \\    \\ \\   __  \\ \\  \\ \\  \\    \\ \\   __  \\ \\  \\        Running: CLI Server" + " by cyt(r3inbowari)")
	Log.Blue("  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\       Listened: " + GetConfig(false).APIAddr)
	Log.Blue("   \\ \\_______\\ \\__\\ \\__\\ \\__\\ \\_______\\ \\__\\ \\__\\ \\__\\      PID: " + strconv.Itoa(os.Getpid()))
	Log.Blue("    \\|_______|\\|__|\\|__|\\|__|\\|_______|\\|__|\\|__|\\|__|      Built on " + Up.BuildTime)
	Log.Blue("")
}

func Shutdown(ctx context.Context) {
	s.Shutdown(ctx)
}
