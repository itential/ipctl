// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/itential/ipctl/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// EnableFileLogs Enables file logging. If torero is being launched as a server, logs will go
// in the directory defined at IAP_LOG_SERVER_DIR. If launched in client mode
// logs will go in IAP_APPLICATION_WORKING_DIR.
func EnableFileLogs(cfg *config.Config) {
	logDir := cfg.WorkingDir

	// Not using internal/functions/EnsureExists to avoid cyclic dependencies
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		fmt.Printf("the logging directory at %s does not exist. Creating the directory now\n", logDir)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "%s. Unable to create logging directory at %s\n", err.Error(), logDir)
			os.Exit(1)
		}
	}

	logFullFilePath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(
		logFullFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s. Failed to open the logging file '%s'\n", err.Error(), logFullFilePath)
		os.Exit(1)
	}

	if cfg.LogFileJSON {
		iowriters = append(iowriters, logFile)
	} else {
		fileWriter := zerolog.NewConsoleWriter()
		fileWriter.Out = logFile
		fileWriter.NoColor = true
		fileWriter.FormatTimestamp = timestampFormatter(cfg.LogTimestampTimezone)

		iowriters = append(iowriters, fileWriter)
	}

	// We only have one iowriter at first but more can be added later
	writers := zerolog.MultiLevelWriter(iowriters...)
	log.Logger = zerolog.New(writers).With().Timestamp().Logger()
}
