package skb

import (
	"log"
	//        "reflect"
	"testing"
	//        "unsafe"
)

func TestRBtree(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	Init_skb_mem_pool()
	skb_head := new(Sk_buff_head)
	skb_head.Queue_init(false)
	newsk1, _ := Alloc(100)
	skb_head.Queue_tail(newsk1)
	newsk2, _ := Alloc(100)
	skb_head.Queue_tail(newsk2)

	log.Printf("--%p---%p----%p\n", skb_head, newsk1, newsk2)
	if skb_head.Queue_empty() {
		log.Println("skb queue is empty")
	} else {
		log.Println("skb queue is not empty")
	}

	skb := skb_head.Dequeue()
	log.Printf("skb:%p\n", skb)

	skb = skb_head.Dequeue()
	log.Printf("skb:%p\n", skb)

	if skb_head.Queue_empty() {
		log.Println("skb queue is empty")
	} else {
		log.Println("skb queue is not empty")
	}

	//  if skb_head == nil {
	//      t.Errorf("Init skb head fail")
	//  }
}
