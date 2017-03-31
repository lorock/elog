## elog
An simple file log system based on golang.

## Installation
go get -u -v github.com/zdt3476/elog

## example
```go
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
	logger.Warn("This is a warn msg.")
	logger.Error("This is a error msg.")
}
```