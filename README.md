# elog
An asynchronous log system based on golang.

## example
```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/zdt3476/elog"
)

func main() {
	dir, _ := os.Getwd()
	path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	filename := filepath.Join(path, "log/log.log")

	cfg := elog.NewDefaultConfig(filename)

	logger, err := elog.NewELog(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("This is a debug msg.")
	logger.Info("This is a info msg.")

    logger.Panic("This is a panic msg.")
	logger.Warn("This is a warn msg.")

	logger.Fatal("This is a fatal msg.")

	select {}
}
```