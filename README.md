## elog
> An simple file log system based on golang.

## Installation
```go
go get -u -v github.com/zdt3476/elog
```

## example
```go
package main

import (
	"log"

	"github.com/zdt3476/elog"
)

func main() {
	logger, err := elog.NewELog(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.Debug("This is a debug msg.")
	logger.Info("This is a info msg.")
	logger.Warn("This is a warn msg.")
	logger.Error("This is a error msg.")
}
```
## Output
```
[DEBG] 2017-03-31 14:24:36.559 /Users/zdt3476/Go/src/github.com/zdt3476/elog/example/basic.go:25. This is a debug msg.
[INFO] 2017-03-31 14:24:36.559 /Users/zdt3476/Go/src/github.com/zdt3476/elog/example/basic.go:26. This is a info msg.
[WARN] 2017-03-31 14:24:36.559 /Users/zdt3476/Go/src/github.com/zdt3476/elog/example/basic.go:27. This is a warn msg.
[EROR] 2017-03-31 14:24:36.559 /Users/zdt3476/Go/src/github.com/zdt3476/elog/example/basic.go:28. This is a error msg.
```