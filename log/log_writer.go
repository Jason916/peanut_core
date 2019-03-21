//jasonxu
package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"sync"
)

const (
	defaultLogFormatWithPrefix = "[%-5s] [%s]:-> %s\n"
)

const (
	LevelInfoMsg              = "INFO"
	LevelTraceMsg             = "TRACE"
	LevelErrorMsg             = "ERROR"
	LevelWarnMsg              = "WARN"
	LevelSuccessMsg           = "SUCC"
	LevelInfo       LevelType = iota
	LevelTrace
	LevelError
	LevelWarn
	LevelSuccess
)

type LevelType uint8

type PLogger struct {
	writer PLogWriter
}

type PLogWriter interface {
	Write(*logEntity) error
	Close() error
}

type logEntity struct {
	msg   string
	time  time.Time
	level LevelType
}

type WriterFile struct {
	*WriterConfig
	logChan            chan *logEntity
	tickChan           *time.Ticker
	logFileUrl         string
	logCurFileUrl      string
	rotate             bool
	rotateFileDateTime time.Time
}

type WriterConfig struct {
	saveInterval   time.Duration
	cacheSize      uint32
	dateFormat     string
	dateTimeFormat string
}

var (
	pool *sync.Pool
)

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &logEntity{}
		},
	}
}

func NewPLogWriterConfig() *WriterConfig {
	return &WriterConfig{
		saveInterval:   time.Second,
		cacheSize:      1024,
		dateFormat:     "2006-01-02",
		dateTimeFormat: "2006-01-02 15:04:05",
	}
}

func getLevelTag(level LevelType) string {
	switch level {
	case LevelInfo:
		return LevelInfoMsg
	case LevelTrace:
		return LevelTraceMsg
	case LevelError:
		return LevelErrorMsg
	case LevelWarn:
		return LevelWarnMsg
	case LevelSuccess:
		return LevelSuccessMsg
	default:
		return ""
	}
}

func (config *WriterConfig) SetSaveInterval(saveInterval time.Duration) {
	config.saveInterval = saveInterval
}

func (config *WriterConfig) SetCacheSize(cacheSize uint32) {
	config.cacheSize = cacheSize
}

func (config *WriterConfig) SetDateFormat(dateFormat string) {
	config.dateFormat = dateFormat
}

func (config *WriterConfig) SetDateTimeFormat(dateTimeFormat string) {
	config.dateTimeFormat = dateTimeFormat
}

func (w *WriterFile) Write(logEntity *logEntity) error {
	select {
	case w.logChan <- logEntity:
		return nil
	default:
		w.writeFile()
		w.Write(logEntity)
		return errors.New("WriterFile log chan is overflow")
	}
}

func (w *WriterFile) Close() error {
	w.tickChan.Stop()
	err := w.writeFile()
	if err != nil {
		fmt.Println("Writer log file failed", err)
	}
	close(w.logChan)
	return err
}

func (w *WriterFile) getDate(t time.Time) string {
	return t.Format(w.dateFormat)
}

func (w *WriterFile) getDateTime(t time.Time) string {
	return t.Format(w.dateTimeFormat)
}

func (w *WriterFile) writeFile() error {
	if len(w.logChan) == 0 {
		return nil
	}
	file, err := os.OpenFile(w.logCurFileUrl, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}

	for len(w.logChan) > 0 {
		select {
		case logEntity, ok := <-w.logChan:
			if !ok {
				return errors.New("WriterFile log chan is closed")
			}
			if w.rotate && logEntity.time.After(w.rotateFileDateTime) {
				w.updateRotateDateTime(logEntity.time)
				file, err = os.OpenFile(w.logCurFileUrl, os.O_APPEND|os.O_CREATE, 0666)
				if err != nil {
					return err
				}
			}
			fMsg := fmt.Sprintf(defaultLogFormatWithPrefix, getLevelTag(logEntity.level), w.getDateTime(logEntity.time), logEntity.msg)
			_, err = file.WriteString(fMsg)
			pool.Put(logEntity)
		default:
			return nil
		}
	}
	return err
}

func (w *WriterFile) updateRotateDateTime(t time.Time) {
	if w.rotate {
		w.rotateFileDateTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(time.Hour*24 - 1)
		w.updateRotateFile(t)
	}
}

func (w *WriterFile) updateRotateFile(t time.Time) {
	if w.rotate {
		w.logCurFileUrl = w.logFileUrl + "_" + w.getDate(t) + ".log"
	}
}

func logPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (w *WriterFile) SetFileDateFormat(dateFormat string) {
	w.WriterConfig.SetDateFormat(dateFormat)
	w.updateRotateFile(w.rotateFileDateTime)
}

func NewPLogWriter(path string, logFileName string, rotate bool, config *WriterConfig) (*WriterFile, error) {
	if !logPathExists(path) {
		return nil, fmt.Errorf("log path not exists: %s", path)
	}
	if config == nil {
		config = NewPLogWriterConfig()
	}
	logFileUrl := filepath.Join(path, logFileName)
	writer := &WriterFile{
		WriterConfig:  config,
		logChan:       make(chan *logEntity, config.cacheSize),
		tickChan:      time.NewTicker(config.saveInterval),
		logFileUrl:    logFileUrl,
		logCurFileUrl: logFileUrl + ".log",
		rotate:        rotate,
	}
	writer.updateRotateDateTime(time.Now())
	go writer.serve()
	return writer, nil
}

func (w *WriterFile) serve() {
	for {
		select {
		case <-w.tickChan.C:
			if err := w.writeFile(); err != nil {
				fmt.Println("Writer log file failed", err)
			}
		}
	}
}

func NewPLogger(plw PLogWriter) *PLogger {
	plogger := &PLogger{
		writer: plw,
	}

	return plogger
}

func (log *PLogger) writeInfo(level LevelType, info string, args ...interface{}) {
	pInfo := fmt.Sprintf(info, args...)
	curtime := time.Now()
	if log.writer != nil {
		entity := pool.Get().(*logEntity)
		entity.msg = pInfo
		entity.level = level
		entity.time = curtime
		if err := log.writer.Write(entity); err != nil {
			fmt.Println("write info failed", err)
		}
	}
}

func (log *PLogger) PInfo(info string, args ...interface{}) {
	log.writeInfo(LevelInfo, info, args...)
}

func (log *PLogger) PTrace(info string, args ...interface{}) {
	log.writeInfo(LevelTrace, info, args...)
}

func (log *PLogger) PError(info string, args ...interface{}) {
	log.writeInfo(LevelError, info, args...)
}

func (log *PLogger) PWarn(info string, args ...interface{}) {
	log.writeInfo(LevelWarn, info, args...)
}

func (log *PLogger) PSucc(info string, args ...interface{}) {
	log.writeInfo(LevelWarn, info, args...)
}

func (log *PLogger) PClose() error {
	if log.writer != nil {
		return log.writer.Close()
	}
	return nil
}
