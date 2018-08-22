// jasonxu-2018/8/21
package log

import (
	"testing"
	"github.com/Jason916/peanut_core/log"
)

func TestLog(t *testing.T) {
	log.Info("test info %s", "test")
	log.Trace("test trace %d",  111)
	log.Warning("test warning")
	log.Success("test success")
}
