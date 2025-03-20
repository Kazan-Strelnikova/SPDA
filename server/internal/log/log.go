package log

import (
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/log/pretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env, host, port string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		var w io.Writer
		conn, err := net.Dial("tcp", host + ":" + port)
		if err != nil {
			w = os.Stdout
		} else {
			w = conn
		}
		log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
		
	case envProd:
		var w io.Writer
		conn, err := net.Dial("tcp", host + ":" + port)
		if err != nil {
			w = os.Stdout
		} else {
			w = conn
		}
		log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := pretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
