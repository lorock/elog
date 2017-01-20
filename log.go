package elog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	goroutineCount = 40000
)

var (
	errNoAbsPath        = errors.New("ELog: Config.AbsPath is empty.")
	errPermissionDenied = errors.New("ELog: Config.Perm is lowest permission, you need change it.")
)

type ELog struct {
	logger *logger

	stdout io.Writer

	logChan chan *logMessage

	cfg *Config

	waitWrite sync.WaitGroup
	waitClose sync.WaitGroup
}

type logger struct {
	f *os.File

	sync.Mutex
}

type logMessage struct {
	msg string

	lvl LogLevel

	lineNo string // 行号信息
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

	dir := filepath.Dir(cfg.AbsPath)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}

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
	e.waitWrite.Add(1)
	go func() {
		defer func() {
			e.waitClose.Done()
		}()

		if e.logChan == nil {
			e.logChan = make(chan *logMessage, goroutineCount)
		}

		e.waitWrite.Done()
		e.waitClose.Add(1)
		for {
			select {
			case logMsg, ok := <-e.logChan:
				if !ok {
					log.Println("Chan is closed.")
					return
				}
				e.log(logMsg)

				// 处理Fatal
				if logMsg.lvl == FatalLvl {
					e.Close()
					os.Exit(-1)
				}

				msgPool.Put(logMsg) // 用完之后将对象放回池里
			}
		}
	}()
}

// func (e *ELog) Wait() {
// 	for len(e.logChan) > 0 {
// 		time.Sleep(time.Millisecond)
// 	}
// }

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
	// e.logger.Lock()
	// defer e.logger.Unlock()

	log.Println("ELog: start close log.")
	defer log.Println("Elog: log closed.")

	var err error

	close(e.logChan)
	// 等待所有输出完成
	e.waitClose.Wait()
	log.Println("After wait.")

	if err = e.logger.f.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Sync log data in Close encounter a error.Error:%v\n", err)
	}

	if err = e.logger.f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Close log file in Close encounter a error.Error:%v\n", err)
	}
}
