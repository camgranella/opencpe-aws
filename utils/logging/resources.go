package logging

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"log/slog"
	"os"
	"strings"
)

// JSON loggers are used for processing files
// The text loggers for actually logging relevant information so users can check in the log output

func JSONInfoLogger() {
	// Pl̀aceholder info
	JSONLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	JSONLogger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"sample message",
		slog.String("method", "GET"),
		slog.Group("properties",
			slog.Int("sample 1", 4000),
			slog.Int("sample 2", 3000),
			slog.String("format", "jpeg"),
		),
	)
}

func TextRequestOutputLogger(msg string, config string, policy string, region string) {
	TextLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	TextLogger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg,
		slog.String("ConfigFile", config),
		slog.String("Policy", policy),
		slog.String("Region", region),
	)
}

func BreakerLine() {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		width = 80 //Fallback terminal width
	}

	line := strings.Repeat("_", width)
	fmt.Println(line)
}
