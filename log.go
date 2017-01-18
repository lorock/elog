package elog

import (
	"errors"
	"io"
	"os"
	"sync"
)

var (
	errNoAbsPath        = errors.New("ELog: Config.AbsPath is empty.")
	errPermissionDenied = errors.New("ELog: Config.Perm is lowest permission, you need change it.")
)

type ELog struct {
	log *logger

	stdout io.Writer

	cfg *Config
}

type logger struct {
	f *os.File

	sync.Mutex
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
		log: log,
		cfg: cfg,
	}

    if cfg.EnabledStdout {
        elog.stdout = os.Stdout
    }

	return elog, nil
}

// reload config.It's aim to rotate log.
func (e *ELog) Reload() error {
	e.log.Lock()
	defer e.log.Unlock()

	var err error

	// flush all data
	if err = e.log.f.Sync(); err != nil {
		return err
	}

	// close file
	if err = e.log.f.Close(); err != nil {
		return err
	}

	// reopen file
	e.log.f, err = os.OpenFile(e.cfg.AbsPath, os.O_WRONLY|os.O_APPEND, e.cfg.Perm)
	if err != nil {
		return err
	}

	return nil
}
	return nil
}
