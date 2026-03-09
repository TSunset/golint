package a

import (
	"log/slog"

	"go.uber.org/zap"
)

type customLogger struct{}

func (c *customLogger) Info(msg string)  {}
func (c *customLogger) Error(msg string) {}

func test(zl *zap.Logger, sl *slog.Logger, custom *customLogger, msg, password, apiKey, token string) {
	slog.Info("starting server")
	slog.Error("Failed to connect to database")                                                  // want "loglint: log message must start with a lowercase letter"
	slog.Info("\u0437\u0430\u043f\u0443\u0441\u043a \u0441\u0435\u0440\u0432\u0435\u0440\u0430") // want "loglint: log message must be in English only"
	slog.Info("server started!!!")                                                               // want "loglint: log message must not contain special characters or emoji"
	slog.Info("token: abc123")                                                                   // want "loglint: log message must not contain sensitive data"

	sl.Info("Another bad Message") // want "loglint: log message must start with a lowercase letter"

	zl.Info("connection failed!!!")         // want "loglint: log message must not contain special characters or emoji"
	zl.Info("user password: secret")        // want "loglint: log message must not contain sensitive data"
	zl.Debug("api_key=" + apiKey)           // want "loglint: log message must not contain sensitive data"
	slog.Info("user password: " + password) // want "loglint: log message must not contain sensitive data"
	slog.Info("token: " + token)            // want "loglint: log message must not contain sensitive data"
	slog.Info("token validated")

	custom.Info("Bad custom logger message")
	custom.Error("token: abc123")

	slog.Info(msg)
	zl.Info(msg)
}
