package benchmark

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zdt3476/elog"
)

func newLog() *elog.ELog {
	dir, _ := os.Getwd()
	path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	filename := filepath.Join(path, "log/log.log")

	cfg := elog.NewDefaultConfig(filename)
	cfg.EnabledStdout = false

	logger, err := elog.NewELog(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func BenchmarkRunParallel(b *testing.B) {
	logger := newLog()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("Go fast.Int(%d),int64(%d),float(%f),string(%s),bool(%v),time(%v),duration(%d)",
				1, int64(1), 3.0, "four!", true, time.Now().UnixNano(), time.Second)
		}
	})
	logger.Close()
}
