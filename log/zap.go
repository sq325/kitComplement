// Copyright 2023 Sun Quan
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func DefaultZapLogger() *zap.SugaredLogger {
	encoder := getEncoder("console")
	writer := os.Stdout
	writeSyncer := zapcore.AddSync(writer)
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	return zap.New(core, zap.AddCaller()).Sugar()
}

// 初始化日志实例
// svcName: 服务名，自动生成服务名.log文件
// logLevel: One of [debug, info, warning, error], default info
// logFormat: One of [json, console], defaul console
// 日志文件默认保存30天，每天切换文件
func NewZapLogger(svcName, logLevel, logFormat string) *zap.SugaredLogger {
	logName := svcName + ".log"
	writer, _ := rotatelogs.New(
		logName+".%Y_%m_%d",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithMaxAge(time.Duration(30)*24*time.Hour),    // 保存logRetentionTimes天
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour), // 日切
	)
	encoder := getEncoder(logFormat)
	writeSyncer := zapcore.AddSync(writer)

	var level = zap.InfoLevel
	switch logLevel {
	case "info":
		level = zap.InfoLevel
	case "warning":
		level = zap.WarnLevel
	case "debug":
		level = zap.DebugLevel
	case "error":
		level = zap.ErrorLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)

	return zap.New(core, zap.AddCaller()).Sugar() // AddCaller() 显示行号和文件名
}

// 日志格式
func getEncoder(logFormat string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式：2020-12-16T17:53:30.466+0800
	// encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder   // 时间格式：2020-12-16T17:53:30.466+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 在日志文件中使用大写字母记录日志级别
	if logFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}
