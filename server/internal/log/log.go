package log

import (
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
		conn, err := net.Dial("tcp", host + ":" + port)
		if err != nil {
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		} else {
			log = slog.New(slog.NewJSONHandler(conn, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}
	case envProd:
		conn, err := net.Dial("tcp", host + ":" + port)
		if err != nil {
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		} else {
			log = slog.New(slog.NewJSONHandler(conn, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}
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
