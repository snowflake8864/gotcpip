package device

import (
	//	"bytes"
	//	"encoding/binary"
	//"fmt"
	"github.com/snowflake8864/gotcpip/constant"
	"hash/crc32"
	//	"tcpip/stack"
	//	"header"
	P "github.com/snowflake8864/libs/public"
	"log"
	//	"tcpip/net/protocol"
	"github.com/snowflake8864/gotcpip/skb"
	"github.com/snowflake8864/libs/rbtree"
	//	"tcpip/utils"
	"time"
	//	"unsafe"
)

const DEVICE_DEFAULT_MTU = 1500

type ProcessFunc func(dev *ZBDevice, _skb *skb.Sk_buff) int

type ZBDevice struct {
	Name     string
	Id       uint8
	Hash     uint32
	Overhead uint32
	Mtu      uint32
	Handle   P.Void
	Ethdev   P.Void /* Null if non-ethernet */
	QueueIn  *skb.Sk_buff_head
	QueueOut *skb.Sk_buff_head

	NetworkQueue   *skb.Sk_buff_head
	TransportQueue *skb.Sk_buff_head

	ChanTransportIn  chan *skb.Sk_buff_head
	ChanTransportOut chan *skb.Sk_buff_head

	ChanNetworkIn  chan *skb.Sk_buff_head
	ChanNetworkOut chan *skb.Sk_buff_head
	LinkState      func(self *ZBDevice) int
	Send           func(self *ZBDevice, buf []byte, len int) int /* Send function. Return 0 if busy */
	Poll           func(self *ZBDevice, loop_score int) int
	ProcessIn      ProcessFunc
	Init           func(dev *ZBDevice, mac []byte) P.Derror
	Destroy        func(self *ZBDevice)
}
type DevicesRRInfo struct {
	nodeIn, nodeOut *rbtree.Node
}

var devicesRRInfo = DevicesRRInfo{
	nodeIn:  nil,
	nodeOut: nil,
}

/*
func deviceInitMac(dev *ZBDevice, mac []byte) int {

dev.Init(dev, mac)
	reader := bytes.NewReader(mac[0:6])
	//	err = binary.Read(reader, binary.BigEndian, &buf)
	err := binary.Read(reader, binary.LittleEndian, &dev.Eth.mac)
	if err != nil {
		log.Fatal(err)
	}
return 0
}
*/
func devCmp(ka, kb rbtree.Item) int {
	var a, b *ZBDevice
	switch v := ka.(type) {
	case *ZBDevice:
		a = v
	default:

		log.Println("unknown")
	}
	switch v := kb.(type) {
	case *ZBDevice:
		b = v
	default:
		log.Println("unknown")
	}

	if a.Hash < b.Hash {
		return -1
	}

	if a.Hash > b.Hash {
		return 1
	}

	return 0
}

var DeviceTree = &rbtree.RBtree{}
var IsInitDeviceTree = false
var devId uint8 = 0

func InitDeviceTree() {
	log.Println("==========InitDeviceTree=========")
	DeviceTree = rbtree.InitTree(devCmp)
	if DeviceTree == nil {
		log.Println("Init DeviceTree tree fail")
	}
	log.Printf("DeviceTree[%v]\n", DeviceTree)
	IsInitDeviceTree = true
}
func DeviceInit(dev *ZBDevice, name string, mac []byte) P.Derror {
	if !IsInitDeviceTree {
		log.Println("Please Init Device Tree!")
		return P.Perror(P.ERR_EEXIST)
	}
	dev.Name = name

	dev.Id = devId
	devId++
	h := crc32.NewIEEE()
	h.Write([]byte(dev.Name))
	dev.Hash = h.Sum32()

	devicesRRInfo.nodeIn = nil
	devicesRRInfo.nodeOut = nil
	dev.QueueIn = new(skb.Sk_buff_head)
	if dev.QueueIn == nil {
		log.Println("dev.QueueIn == nil")
		return P.Perror(P.ERR_EEXIST)
	}
	dev.QueueIn.Queue_init(false)
	dev.ChanNetworkIn = make(chan *skb.Sk_buff_head, 1)
	dev.ChanNetworkOut = make(chan *skb.Sk_buff_head, 1)
	dev.ChanTransportIn = make(chan *skb.Sk_buff_head, 1)
	dev.ChanTransportOut = make(chan *skb.Sk_buff_head, 1)

	log.Printf("dev.ChanIn:%p\n", dev.ChanNetworkIn)
	dev.QueueOut = new(skb.Sk_buff_head)
	if dev.QueueOut == nil {
		log.Println("dev.QueueOut == nil")
		return P.Perror(P.ERR_EEXIST)
	}
	dev.QueueOut.Queue_init(false)
	/*
		DeviceTree = rbtree.InitTree(devCmp)
		if DeviceTree == nil {
			log.Println("Init DeviceTree tree fail")
		}
		log.Printf("DeviceTree[%v]--dev[%v]\n", DeviceTree, dev)
	*/
	DeviceTree.Insert(dev)
	if dev.Mtu == 0 {
		dev.Mtu = DEVICE_DEFAULT_MTU
	}

	if mac != nil && dev.Init != nil && dev.Init(dev, mac) == nil {
		return P.Perror(P.ERR_EEXIST)
		//		ret = deviceInitMac(dev, mac)
	} else {
		log.Println("mac == nil")
		dev.Ethdev = nil
		return nil
		//return P.Perror(P.ERR_EEXIST)
		//ret = deviceInitNomac(dev)
	}

	return nil
}

/*
func deviceDestroy(dev *ZBDevice) {

	SkbQueuePurge(dev.QueueIn)
	SkbQueuePurge(dev.QueueOut)

	ip.Ipv4CleanupLinks(dev)

	Device_tree.rbtree.Delete(dev)

	if dev.destroy != nil {
		dev.destroy(dev)
	}

	DevicesRRInfo.nodeIn = nil
	DevicesRRInfo.nodeOut = nil
	// FREE(dev);
}
*/
const DEV_LOOP_MIN = 32

func DevicesLoop(loopScore, direction int) int {

	startLoopScore := loopScore
	nextNode := devRoundrobinStart(direction)
	//	startLoopScore := loopScore
	if nextNode == nil {
		return loopScore
	}
	next := nextNode.KeyValue
	start := next
	loopScore = 64
	/* round-robin all devices, break if traversed all devices */
	for loopScore > DEV_LOOP_MIN && next != nil {
		loopScore = devLoop(next.(*ZBDevice), loopScore, direction)
		nextNode = rbtree.TreeNext(nextNode)
		next = nextNode.KeyValue
		if next == nil {
			nextNode = rbtree.TreeFirstNode(DeviceTree.Root)
			next = nextNode.KeyValue
		}

		if next == start {
			loopScore = startLoopScore
			break
		}
	}
	devRoundrobinEnd(direction, nextNode)
	return loopScore
}

func devRoundrobinStart(direction int) *rbtree.Node {

	if devicesRRInfo.nodeIn == nil {
		devicesRRInfo.nodeIn = rbtree.TreeFirstNode(DeviceTree.Root)
	}

	if devicesRRInfo.nodeOut == nil {
		devicesRRInfo.nodeOut = rbtree.TreeFirstNode(DeviceTree.Root)
	}

	if direction == constant.LOOP_DIR_IN {
		return devicesRRInfo.nodeIn
	} else {
		return devicesRRInfo.nodeOut
	}
}

func devRoundrobinEnd(direction int, last *rbtree.Node) {
	if direction == constant.LOOP_DIR_IN {
		devicesRRInfo.nodeIn = last
	} else {
		devicesRRInfo.nodeOut = last
	}
}

func devLoop(dev *ZBDevice, loopScore int, direction int) int {
	if dev.Poll != nil {
		//loopScore = dev.Poll(dev, loopScore)
		go dev.Poll(dev, loopScore)
	} else {
		log.Printf("No register poll for this device\n")
		return -1
	}
	if direction == constant.LOOP_DIR_OUT {
		//loopScore = devloopOut(dev, loopScore)
		go devloopOut(dev, loopScore)
	} else {
		//loopScore = devloopIn(dev, loopScore)
		go devloopIn(dev, loopScore)
	}

	return loopScore
}

func devloopIn(dev *ZBDevice, loopScore int) int {
	//dev.NetworkQueue = <-dev.ChanNetworkIn
	//log.Printf("---------dev.QueueUp:%p\n--------", dev.NetworkQueue)
	/*	var QueueUp interface{}
		QueueUp = <-dev.ChanIn
		log.Println("type:%v", QueueUp.(type))
	*/
	if dev.QueueIn == nil {
		return loopScore
	}
	for loopScore > 0 {
		if dev.QueueIn.Queue_empty() == true {
			//break
			time.Sleep(1 * time.Second)
			continue
		}

		//log.Println("dequeu dev.QueueIn")
		_skb := dev.QueueIn.Dequeue()
		if _skb != nil && dev.ProcessIn != nil {
			loopScore--
			//_skb.Datalink_header = _skb.Data
			dev.ProcessIn(dev, _skb)
			log.Println(_skb.Data)
			//stack.EthernetReceive(_skb)
			//loopScore--
		}
	}
	return loopScore
}

func devloopOut(dev *ZBDevice, loopScore int) int {

	if dev.QueueOut == nil {
		return loopScore
	}

	for loopScore > 0 {
		if dev.QueueOut.Queue_empty() == true {
			//break
			time.Sleep(1 * time.Second)
			continue
		}
		_skb := dev.QueueOut.Peek()
		if _skb != nil && dev.Send(dev, _skb.Data, int(_skb.Data_len)) == 0 { /* success. */
			_skb := dev.QueueOut.Dequeue()
			_skb.Free() /* SINGLE POINT OF DISCARD for OUTGOING FRAMES */
			loopScore--
		} else {
			break
		} /* Don't discard */

	}
	return loopScore
}

/*
func pico_device *pico_get_device(const char*name) *ZBDevice {

    struct pico_device *dev;
    struct pico_tree_node *index;


	pos := rbtree.TreeFirstNode(TreeDevLink.Root)
	for pos != &LEAF {
		link := index.KeyValue

		if dev == link.dev {

			ipv4LinkDel(dev, link.address)

		}
		pos = rbtree.TreeNext(pos)
	}



    pico_tree_foreach(index, &Device_tree){
        dev = index->keyValue;
        if(strcmp(name, dev->name) == 0)
            return dev;
    }
    return NULL;
}
*/
/*
type SKBQueue struct {
	C chan *skb.Sk_buff
}
*/
