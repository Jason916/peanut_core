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
	defaultCacheSize           = 1024
	defaultSaveInterval        = time.Second * 1
	defaultLogFormatPrefixFile = "[%-5s] [%s] : %s -> %s \n"
	defaultDateFormat          = "2010-01-01"
	defaultDateTimeFormat      = "2010-01-01 12:00:00.000"
)

var (
	pool *sync.Pool
)

type LevelType uint8

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

type logEntity struct {
	msg    string
	time   time.Time
	level  LevelType
	caller string
}

type WriterFile struct {
	*WriterConfig
	logChan        chan *logEntity
	tickChan       *time.Ticker
	level          LevelType
	fileUrl        string
	curFileUrl     string
	rotate         bool
	rotateFileDate time.Time
}

type WriterConfig struct {
	cacheSize      uint32
	saveInterval   time.Duration
	dateFormat     string
	dateTimeFormat string
}

func NewLogWriterConfig() *WriterConfig {
	return &WriterConfig{
		saveInterval:   defaultSaveInterval,
		cacheSize:      defaultCacheSize,
		dateFormat:     defaultDateFormat,
		dateTimeFormat: defaultDateTimeFormat,
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

func (w *WriterFile) Write(logEntity *logEntity) error {
	if logEntity.level < w.level {
		return nil
	}
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
	close(w.logChan)
	if err != nil {
		fmt.Println("Writer log file failed", err)
	}
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
	file, err := os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for len(w.logChan) > 0 {
		select {
		case logEntity, ok := <-w.logChan:
			if !ok {
				return errors.New("WriterFile log chan is closed")
			}
			if w.rotate && logEntity.time.After(w.rotateFileDate) {
				w.refreshRotateDate(logEntity.time)
				file, err = os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE, 0666)
				if err != nil {
					return err
				}
			}
			fMsg := fmt.Sprintf(defaultLogFormatPrefixFile, getLevelTag(logEntity.level), w.getDateTime(logEntity.time), logEntity.caller, logEntity.msg)
			_, err = file.WriteString(fMsg)
			pool.Put(logEntity)
		default:
			return nil
		}
	}
	return err
}

func (w *WriterFile) refreshRotateDate(t time.Time) {
	if w.rotate {
		w.rotateFileDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(time.Hour*24 - 1)
		w.refreshRotateFile(t)
	}
}

func (w *WriterFile) refreshRotateFile(t time.Time) {
	if w.rotate {
		w.curFileUrl = w.fileUrl + "_" + w.getDate(t) + ".log"
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (w *WriterFile) SetDateFormat(dateFormat string) {
	w.WriterConfig.SetDateFormat(dateFormat)
	w.refreshRotateFile(w.rotateFileDate)
}

func NewLogWriterFile(level LevelType, path string, fileName string, rotate bool, config *WriterConfig) (*WriterFile, error) {
	if !pathExists(path) {
		return nil, fmt.Errorf("path not exists: %s", path)
	}
	if config == nil {
		config = NewLogWriterConfig()
	}
	fileUrl := filepath.Join(path, fileName)
	writer := &WriterFile{
		WriterConfig: config,
		logChan:      make(chan *logEntity, config.cacheSize),
		tickChan:     time.NewTicker(config.saveInterval),
		level:        level,
		fileUrl:      fileUrl,
		curFileUrl:   fileUrl + ".log",
		rotate:       rotate,
	}
	writer.refreshRotateDate(time.Now())
	go writer.serve()
	return writer, nil
}
