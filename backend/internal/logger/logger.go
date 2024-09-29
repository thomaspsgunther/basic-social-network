package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
	"y_net/internal/utils"
)

var (
	ServerLogger       *CustomLogger
	logFile            *os.File
	logDir             string
	currentLogFileName string
	mu                 sync.Mutex // Mutex for synchronization
)

type CustomLogger struct {
	*slog.Logger
}

// NewCustomLogger initializes a new custom logger
func NewCustomLogger(logFile io.Writer) *CustomLogger {
	return &CustomLogger{
		Logger: slog.New(slog.NewTextHandler(logFile, nil)),
	}
}

// Fatalf formats and logs a fatal message, then exits
func (cl *CustomLogger) Fatalf(format string, args ...interface{}) {
	cl.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Fatal logs a fatal message without formatting, then exits
func (cl *CustomLogger) Fatal(args ...interface{}) {
	cl.Error(fmt.Sprint(args...))
	os.Exit(1)
}

// InitLogger initializes the logger with rotation options
func InitLogger(rotationFrequency string) error {
	var err error
	logDir, err = findLogDirectory()
	if err != nil {
		return err
	}

	// Prepare the initial log file name with date, hour, and minute
	initialLogFileName := fmt.Sprintf("server-%s.log", time.Now().Format("2006-01-02_15-04"))
	currentLogFileName = initialLogFileName
	initialLogFilePath := filepath.Join(logDir, initialLogFileName)

	mu.Lock() // Locking before file operations
	logFile, err = os.OpenFile(initialLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	mu.Unlock() // Unlocking after file operations
	if err != nil {
		return err
	}

	ServerLogger = NewCustomLogger(io.MultiWriter(logFile, os.Stderr))

	ServerLogger.Info("--------------------------------------------------------------------")
	ServerLogger.Info("logger initialized")

	// Start log rotation in a separate goroutine
	go rotateLogs(rotationFrequency)

	return nil
}

// findLogDirectory finds the log directory
func findLogDirectory() (string, error) {
	rootPath, err := utils.FindProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(rootPath, "logs"), nil
}

// rotateLogs checks the current time and rotates logs as necessary
func rotateLogs(frequency string) {
	var duration time.Duration
	switch frequency {
	case "hourly":
		duration = time.Hour
	case "daily":
		duration = 24 * time.Hour
	case "weekly":
		duration = 7 * 24 * time.Hour
	case "monthly":
		duration = 30 * 24 * time.Hour // Approximation for monthly rotation
	default:
		duration = 24 * time.Hour // Default to daily
	}

	for {
		time.Sleep(duration)
		err := rotateLogFile()
		if err != nil {
			ServerLogger.Error("Error rotating log file: " + err.Error())
		}
	}
}

// rotateLogFile handles the log rotation and gzip compression
func rotateLogFile() error {
	mu.Lock()         // Locking for file operations
	defer mu.Unlock() // Ensure unlocking happens after the function completes

	// Close the current log file
	if logFile != nil {
		ServerLogger.Info("--------------------------------------------------------------------")
		ServerLogger.Info("rotating log file")

		logFile.Close()
	}

	// Prepare the new log file name with the current timestamp
	newLogFileName := fmt.Sprintf("server-%s.log", time.Now().Format("2006-01-02_15-04"))
	newLogFilePath := filepath.Join(logDir, newLogFileName)

	// Create a new log file with the timestamp
	var err error
	logFile, err = os.OpenFile(newLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Gzip the current log file if it exists
	oldLogFilePath := filepath.Join(logDir, currentLogFileName)
	if _, err := os.Stat(oldLogFilePath); err == nil {
		// Prepare gzipped file name with the old log file's timestamp
		gzippedLogFileName := fmt.Sprintf("%s.gz", oldLogFilePath)
		if err := gzipLogFile(oldLogFilePath, gzippedLogFileName); err != nil {
			return err
		}
	}

	// Delete the old log file after gzipping
	if err := os.Remove(oldLogFilePath); err != nil {
		return err
	}

	// Update currentLogFileName to the new log file
	currentLogFileName = newLogFileName

	ServerLogger = NewCustomLogger(io.MultiWriter(logFile, os.Stderr))
	return nil
}

// gzipLogFile compresses the old log file using gzip
func gzipLogFile(logFilePath string, gzippedLogFilePath string) error {
	inFile, err := os.Open(logFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(gzippedLogFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := gzip.NewWriter(outFile)
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}

// CloseLogger closes the logger
func CloseLogger() {
	mu.Lock()         // Locking before closing
	defer mu.Unlock() // Ensure unlocking after the function completes
	if logFile != nil {
		logFile.Close()
	}
}
