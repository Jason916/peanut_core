//jasonxu
package log

import (
	"testing"
	"fmt"
)

func TestLogWriter(t *testing.T){
	plogWriteConfig := NewPLogWriterConfig()
	plogWriter, err := NewPLogWriter(".", "test_log", true, plogWriteConfig)
	if err != nil{
		t.Error("write log failed, cause by: ", err)
		return
	}
	plog := NewPLogger(plogWriter)
	plog.PInfo("test info ")
	plog.PError("test write error")
	plog.PTrace("test write trace")
	if err := plog.PClose(); err != nil {
		fmt.Println("log close failed, cause by: ", err)
	}
	fmt.Println("test finished")
}
