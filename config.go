package elog

import "os"

type Config struct {
	LogLevel LogLevel // 日志等级，文件和标准输出共用

	AbsPath string // 日志文件的绝对路径

	Perm os.FileMode // 文件权限

    EnabledStdout bool // 是否在控制台输出
}

func NewDefaultConfig(absPath string) *Config {
	return &Config{
		LogLevel: DebugLvl,
		AbsPath:  absPath,
		Perm:     0644,
	}
}
