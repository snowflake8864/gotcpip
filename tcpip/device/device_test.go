package device

import (
	//	"flag"
	"github/snowflake8864/libs/rbtree"
	"log"
	"time"
	//"tcpip/net/ethernet"
	"github/snowflake8864/gotcpip/constant"
	"testing"
	//p "public"
)

func TestRecive(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rbtree.InitRBtreeMemPool()
	dev := new(ZBDevice)
	mac := []byte{00, 0x1c, 0x42, 0xe3, 0x6b, 0x8d}
	err := DeviceInit(dev, "eth0", mac)
	if err != nil {
		t.Errorf("device>  init device failed")
	}
	DevicesLoop(100, constant.LOOP_DIR_IN)
	for {
		time.Sleep(10 * time.Minute)
	}
}
