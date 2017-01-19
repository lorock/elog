package elog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	goroutineCount = 6000
)

var (
	errNoAbsPath        = errors.New("ELog: Config.AbsPath is empty.")
	errPermissionDenied = errors.New("ELog: Config.Perm is lowest permission, you need change it.")
)

type ELog struct {
	logger *logger

	stdout io.Writer

	logChan chan logMessage

	cfg *Config
}

type logger struct {
	f *os.File

	sync.Mutex
}

type logMessage struct {
	msg string

	lvl LogLevel
}

func NewELog(cfg *Config) (*ELog, error) {
	if cfg == nil {
		cfg = NewDefaultConfig("/var/log/elog.log")
	}

	if len(cfg.AbsPath) < 1 {
		return nil, errNoAbsPath
	}

	if cfg.Perm <= 0 {
		return nil, errPermissionDenied
	}

	var (
		fileLog *os.File
		err     error
		elog    *ELog
	)

	fileLog, err = os.OpenFile(cfg.AbsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, cfg.Perm)
	if err != nil {
		return nil, err
	}

	log := &logger{
		f: fileLog,
	}

	elog = &ELog{
		logger: log,
		cfg:    cfg,
	}

	if cfg.EnabledStdout {
		elog.stdout = os.Stdout
	}

	elog.startGoroutine()

	return elog, nil
}

func (e *ELog) startGoroutine() {
	e.logChan = make(chan logMessage, goroutineCount)

	for {
		select {
		case logMsg, ok := <-e.logChan:
			if !ok {
				return
			}
			e.log(logMsg)
			msgPool.Put(logMsg) // 用完之后将对象放回池里
		}
	}
}

// reload config.It's aim to rotate log.
//
// example:
//     sig := make(chan os.Signal)
//     signal.Notify(sig, syscall.SIGHUP)
//     for s := range sig {
//         if s = syscall.SIGHUP {
//             if err := elog.Reload(); err != nil {
//                 fmt.Printf("Rotate log encounter a error.Error: %v", err)
//             }
//         }
//     }
func (e *ELog) Reload() error {
	e.logger.Lock()
	defer e.logger.Unlock()

	var err error

	// flush all data
	if err = e.logger.f.Sync(); err != nil {
		return err
	}

	// close file
	if err = e.logger.f.Close(); err != nil {
		return err
	}

	// reopen file
	e.logger.f, err = os.OpenFile(e.cfg.AbsPath, os.O_WRONLY|os.O_APPEND, e.cfg.Perm)
	if err != nil {
		return err
	}

	return nil
}

func (e *ELog) Close() {
	e.logger.Lock()
	defer e.logger.Unlock()

	fmt.Println("ELog: start close log.")
	defer fmt.Println("Elog: log closed.")

	var err error

	close(e.logChan)

	if err = e.logger.f.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Sync log data in Close encounter a error.Error:%v\n", err)
	}

	if err = e.logger.f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Close log file in Close encounter a error.Error:%v\n", err)
	}
}
