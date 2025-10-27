package main

import (
	"context"
	"cutterproject/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	zapLog "go.uber.org/zap"
)

func main() {
	// Flush zap buffered log first then cancel the context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fiber := config.NewFiber()
	zap := config.NewZap()
	koanf := config.NewKoanf(zap)
	rds := config.NewRedisClient(koanf, zap)
	postgresql := config.NewPostgresqlPool(koanf, zap)

	fiber.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// 5. Compression middleware (should be before logging)
	fiber.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	config.Server(&config.ServerConfig{
		Router:  fiber,
		DB:      postgresql,
		DBCache: rds,
		Log:     zap,
		Config:  koanf,
	})

	GO_SERVER_PORT := koanf.String("GO_SERVER")

	var err error
	go func() {
		err = fiber.Listen(GO_SERVER_PORT)
		if err != nil {
			zap.Fatal("Error Starting Server", zapLog.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	zap.Info("Got one of stop signals")

	err = fiber.ShutdownWithContext(ctx)
	if err != nil {
		zap.Warn("Timeout, forced kill!", zapLog.Error(err))
		zap.Sync()
		os.Exit(1)
	}

	zap.Info("Server has shut down gracefully")
	zap.Sync()
}
