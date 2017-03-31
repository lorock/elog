## ELOG

> 基于Golang的简单易用的文件日志系统。

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

## Config

```go
type Config struct {
	LogLevel LogLevel // 日志等级，文件和标准输出共用 default:DebugLvl

	AbsPath string // 日志文件的绝对路径

	Perm os.FileMode // 文件权限 default:0644

	EnabledStdout bool // 是否在控制台输出 default:true

	ShowLineNumber bool // 显示行号 default:true

	ShortFileName bool // 显示行号时，文件名是否包含路径 default:false

	TimeLayout string // 控制时间显示格式，default:2006-01-02 15:04:05.999
}
```

## Rotate

> 推荐使用logrotate工具进行文件切分,使用Elog.Reload()方法重定向输出到新文件即可。
```go
sig := make(chan os.Signal)
signal.Notify(sig, syscall.SIGHUP)
for s := range sig {
	if s = syscall.SIGHUP {
		if err := elog.Reload(); err != nil {
			fmt.Fprintf(os.Stderr, "Rotate log encounter a error.Error: %v", err)
		}
	}
}
```