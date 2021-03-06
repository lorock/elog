// Copyright (C) 2017, No.20 <zdt3476@163.com>. All rights reserved.

package elog

import "os"

type Config struct {
	LogLevel LogLevel // 日志等级，文件和标准输出共用 default:DebugLvl

	AbsPath string // 日志文件的绝对路径

	Perm os.FileMode // 文件权限 default:0644

	EnabledStdout bool // 是否在控制台输出 default:true

	ShowLineNumber bool // 显示行号 default:true

	ShortFileName bool // 显示行号时，文件名是否包含路径 default:false

	TimeLayout string // 控制时间显示格式，default:2006-01-02 15:04:05.999
}

func NewDefaultConfig(absPath string) *Config {
	return &Config{
		LogLevel:       DebugLvl,
		AbsPath:        absPath,
		Perm:           0644,
		EnabledStdout:  true,
		ShowLineNumber: true,
	}
}
