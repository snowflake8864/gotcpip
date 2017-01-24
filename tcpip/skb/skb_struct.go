package skb

import (
	P "github.com/snowflake8864/libs/public"
	"reflect"
	"sync"
	"unsafe"
)

const (
	SKB_FLAG_BCAST             = (0x01)
	SKB_FLAG_EXT_BUFFER        = (0x02)
	SKB_FLAG_EXT_USAGE_COUNTER = (0x04)
	SKB_FLAG_SACKED            = (0x80)
)
const SBK_STRUCT_SIZE = unsafe.Sizeof(Sk_buff{})

type Info struct {
	SendTTL   uint8
	SendTos   uint8
	LocalAddr uint32
	//	SrcMac    [6]byte
	//	DstMac    [6]byte
	//	ProtoMac  [2]byte
}

func byte2Pointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	)
}

func (this *Sk_buff) GET_INFO() *Info {
	return (*Info)(byte2Pointer(this.CB[0:]))
}

func (this *Sk_buff) SET_TTL(ttl uint8) {
	(*Info)(byte2Pointer(this.CB[0:])).SendTTL = ttl
}

/*
func (this *Sk_buff) GET_SMAC() []byte {
	return (*Info)(byte2Pointer(this.CB[0:])).SrcMac[0:]
}

func (this *Sk_buff) SET_SMAC() {
		SrcMac := (*Info)(byte2Pointer(this.CB[0:])).SrcMac
		SrcMac[0], SrcMac[1], SrcMac[2], SrcMac[3], SrcMac[4], SrcMac[5] =
			this.Data[6], this.Data[7], this.Data[8], this.Data[9], this.Data[10], this.Data[11]
	//((*Info)(byte2Pointer(this.CB[0:])).SrcMac)[0], ((*Info)(byte2Pointer(this.CB[0:])).SrcMac)[1] =
	//	this.Data[6], this.Data[7]
	log.Println("SET_SMAC:", (*Info)(byte2Pointer(this.CB[0:])).SrcMac)
}

func (this *Sk_buff) GET_DMAC() []byte {
	return (*Info)(byte2Pointer(this.CB[0:])).DstMac
}
func (this *Sk_buff) SET_DMAC(addr []byte) {
	(*Info)(byte2Pointer(this.CB[0:6])).DstMac = addr
}
*/
func (this *Sk_buff) GET_TTL() uint8 {
	return (*Info)(byte2Pointer(this.CB[0:])).SendTTL
}
func (this *Sk_buff) SET_LocalAddr(addr uint32) {
	(*Info)(byte2Pointer(this.CB[0:])).LocalAddr = addr
}
func (this *Sk_buff) GET_LocalAddr() uint32 {
	return (*Info)(byte2Pointer(this.CB[0:])).LocalAddr
}

func (this *Sk_buff) SET_TOS(tos uint8) {
	(*Info)(byte2Pointer(this.CB[0:])).SendTos = tos
}

func (this *Sk_buff) GET_TOS() uint8 {
	return (*Info)(byte2Pointer(this.CB[0:])).SendTos
}

//type sk_buff_data_t *uint8
type sk_buff_data_t P.Void
type Sk_buff struct {
	next      *Sk_buff
	prev      *Sk_buff
	Sock      P.Void // *Socket
	tstamp    P.Dtime
	Timestamp uint64
	Dev       P.Void //*device
	CB        [48]uint8
	len,
	Data_len uint32
	mac_len,
	hdr_len uint16
	Frag     uint16
	csum     uint32
	priority uint32
	local_df,
	cloned,
	nohdr bool
	Protocol    uint16
	Payload_len uint16
	App_header,
	Payload,
	TcpOption,
	Transport_header,
	Network_header,
	Datalink_header,
	/* These elements must be at the end, see alloc_skb() for details.  */
	tail,
	end sk_buff_data_t
	NetworkLen,
	TransportLen uint16
	head,
	Data []byte
	addr          []byte
	truesize      uint32
	users         int32 //atomic_t
	Flags         uint8
	CheckPotioner P.Void
}

type Sk_buff_head struct {
	/* These two members must be first. */
	next    *Sk_buff
	prev    *Sk_buff
	shared  bool
	MaxSize uint32
	qlen    uint32
	mutex   sync.Mutex
}
