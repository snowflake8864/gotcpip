package skb

import (
	//	"fmt"
	slab "github.com/couchbase/go-slab"
	"github.com/snowflake8864/gotcpip/utils"
	p "github.com/snowflake8864/libs/public"
	"log"
	//"sync"
	"sync/atomic"
	"unsafe"
)

var skb_head_pool, skb_data_pool *slab.Arena

func (list_ *Sk_buff_head) Queue_len() uint32 {
	return list_.qlen
}

func (list *Sk_buff_head) Queue_init(needlock bool) {
	list.shared = needlock
	list.prev = (*Sk_buff)(unsafe.Pointer(list))
	list.next = (*Sk_buff)(unsafe.Pointer(list))
	//list.next = nil
	//list.prev = nil
	//fmt.Printf("%p-%p-%p\n", list, list.next, list.prev)
	list.qlen = 0
}

func (list *Sk_buff_head) insert(newsk,
	prev, next *Sk_buff) {
	newsk.next = next
	newsk.prev = prev
	next.prev = newsk
	prev.next = newsk
	list.qlen++
}

func (list *Sk_buff_head) Queue_before(
	next,
	newsk *Sk_buff) {
	list.insert(newsk, next.prev, next)
}

func (list_ *Sk_buff_head) Peek() *Sk_buff {
	list := ((*Sk_buff)(unsafe.Pointer(list_))).next
	if list == (*Sk_buff)(unsafe.Pointer(list_)) {
		list = nil
	}
	return list
}
func (list *Sk_buff_head) unlink(skb *Sk_buff) {
	var next, prev *Sk_buff

	list.qlen--
	next = skb.next
	prev = skb.prev
	skb.next = nil
	skb.prev = nil
	next.prev = prev
	prev.next = next
}

func (list *Sk_buff_head) dequeue() *Sk_buff {
	skb := list.Peek()
	if skb != nil {
		list.unlink(skb)
	}
	return skb
}

func (list *Sk_buff_head) queue_tail(newsk *Sk_buff) {
	list.Queue_before((*Sk_buff)(unsafe.Pointer(list)), newsk)
}

func (list *Sk_buff_head) Dequeue() *Sk_buff {
	var skb *Sk_buff

	if list.shared {
		list.mutex.Lock()
		skb = list.dequeue()
		list.mutex.Unlock()
	} else {

		skb = list.dequeue()
	}
	return skb
}
func (list *Sk_buff_head) Queue_tail(newsk *Sk_buff) {

	if list.shared {
		list.mutex.Lock()
		list.queue_tail(newsk)
		list.mutex.Unlock()
	} else {
		list.queue_tail(newsk)
	}
}
func (list *Sk_buff_head) Queue_purge() {
	for {
		skb := list.Dequeue()
		if skb == nil {
			break
		}
		skb.Free()
	}
}
func (list *Sk_buff_head) Queue_empty() bool {
	//fmt.Printf("skb---%p--%p---\n", list.next, list)
	return list.next == (*Sk_buff)(unsafe.Pointer(list))
}

func (list *Sk_buff_head) queue_splice(prev,
	next *Sk_buff) {
	first := list.next
	last := list.prev

	first.prev = prev
	prev.next = first

	last.next = next
	next.prev = last
}

func Skb_queue_splice(list *Sk_buff_head,
	head *Sk_buff_head) {
	if list.Queue_empty() == false {
		list.queue_splice((*Sk_buff)(unsafe.Pointer(head)), head.next)
		head.qlen += list.qlen
	}
}

//type Derror p.Derror
func Alloc(size uint32) (*Sk_buff, p.Derror) {

	/* Get the HEAD */
	buf := skb_head_pool.Alloc(int(SBK_STRUCT_SIZE))
	if buf == nil {
		log.Println("alloc skb struct fail")
		return nil, p.Perror(p.ERR_ENOMEM)
	}

	skb := (*Sk_buff)(utils.Byte2Pointer(buf))
	skb.addr = buf
	data := skb_data_pool.Alloc(int(size))
	if data == nil {
		log.Println("alloc skb data fail")
		return nil, p.Perror(p.ERR_ENOMEM)
	}
	skb.truesize = size + uint32(SBK_STRUCT_SIZE)
	atomic.AddInt32(&skb.users, 1)
	skb.head = data
	skb.Data = data
	skb.Data_len = size
	skb.tail = skb.Data
	skb.end = skb.tail.([]byte)[size:]
	return skb, nil
}

func (skb *Sk_buff) Copy() (*Sk_buff, p.Derror) {
	/* Get the HEAD */
	buf := skb_head_pool.Alloc(int(SBK_STRUCT_SIZE))
	if buf == nil {
		log.Println("alloc skb struct fail")
		return nil, p.Perror(p.ERR_ENOMEM)
	}
	nskb := (*Sk_buff)(utils.Byte2Pointer(buf))
	nskb.addr = buf
	nskb.truesize = skb.truesize
	atomic.AddInt32(&nskb.users, 1)
	nskb.Sock = skb.Sock
	nskb.head = skb.head
	nskb.Data = skb.Data
	nskb.Data_len = skb.Data_len
	nskb.tail = skb.Data
	nskb.end = skb.end
	nskb.App_header = skb.App_header
	nskb.Payload = skb.Payload
	nskb.TcpOption = skb.TcpOption
	nskb.Transport_header = skb.Transport_header
	nskb.Network_header = skb.Network_header
	nskb.Datalink_header = skb.Datalink_header
	nskb.TransportLen = skb.TransportLen
	nskb.NetworkLen = skb.NetworkLen
	nskb.SET_TTL(skb.GET_TTL())
	nskb.SET_TOS(skb.GET_TOS())
	return nskb, nil
}

func (skb *Sk_buff) Free() {

	if skb == nil {
		return
	}
	v := atomic.LoadInt32(&skb.users)
	if v == 1 {
		skb_data_pool.DecRef(skb.Data)
		skb_head_pool.DecRef(skb.addr)
		return
	} else {
		atomic.AddInt32(&skb.users, -1)
		v = atomic.LoadInt32(&skb.users)
		if v != 0 {
			return
		}
	}
	skb_data_pool.DecRef(skb.Data)
	skb_head_pool.DecRef(skb.addr)
}

func Init_skb_mem_pool() {
	skb_head_pool = slab.NewArena(
		int(SBK_STRUCT_SIZE),        // The smallest chunk size is 64B.
		int(SBK_STRUCT_SIZE)*1024*1, // The largest chunk size is 64KB.
		2, // Power of 2 growth in chunk size.
		nil,
	)
	skb_data_pool = slab.NewArena(
		128,        // The smallest chunk size is 64B.
		128*1024*2, // The largest chunk size is 64KB.
		2,          // Power of 2 growth in chunk size.
		nil,
	)

}
