package a

import (
	"log/slog"

	"go.uber.org/zap"
)

type customLogger struct{}

func (c *customLogger) Info(msg string)  {}
func (c *customLogger) Error(msg string) {}
func (c *customLogger) Warn(msg string)  {}
func (c *customLogger) Debug(msg string) {}

func test(zl *zap.Logger, sl *slog.Logger, custom *customLogger, msg, token, apiKey string) {
	slog.Info("pkg slog")      // want "matched slog"
	slog.Error("pkg slog err") // want "matched slog"

	sl.Info("instance slog")   // want "matched slog"
	sl.Warn("instance slog 2") // want "matched slog"

	zl.Info("zap info")           // want "matched zap"
	zl.Error("zap error")         // want "matched zap"
	zl.Debug("api_key=" + apiKey) // want "matched zap"
	slog.Info("token: " + token)  // want "matched slog"

	custom.Info("custom info")
	custom.Error("custom error")

	slog.Info(msg)
	zl.Info(msg)
}
