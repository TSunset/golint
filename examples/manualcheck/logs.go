//go:build manualcheck
// +build manualcheck

package manualcheck

import "log/slog"

func demo(password, token string) {
	slog.Info("service started")
	slog.Info("Starting server on port 8080")
	slog.Error("ошибка подключения к базе данных")
	slog.Warn("connection failed!!!")
	slog.Info("user password: " + password)
	slog.Info("token: " + token)
	slog.Info("token validated")
	slog.Info("Service stopped")
}
