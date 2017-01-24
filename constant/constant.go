package constant

const LOOP_DIR_IN = 1
const LOOP_DIR_OUT = 2
const (
	ZB_IPV4_INADDR_ANY = uint32(0x00000000)
	ZB_IP4_ANY         = 0x00000000
	ZB_IPV4_MTU        = uint16(1500)

	ZB_SIZE_ETHHDR = 14
	//const ZB_SIZE_IP4HDR = uint16(len(IPHeader))
	ZB_SIZE_IP4HDR      = uint16(20)
	ZB_IPV4_MAXPAYLOAD  = uint32(ZB_IPV4_MTU - ZB_SIZE_IP4HDR)
	ZB_IPV4_DONTFRAG    = 0x4000
	ZB_IPV4_MOREFRAG    = 0x2000
	ZB_IPV4_EVIL        = 0x8000
	ZB_IPV4_FRAG_MASK   = 0x1FFF
	ZB_IPV4_DEFAULT_TTL = 64

	/*InetSkbParm flag*/
	IPSKB_FORWARDED        = 1
	IPSKB_XFRM_TUNNEL_SIZE = 2
	IPSKB_XFRM_TRANSFORMED = 4
	IPSKB_FRAG_COMPLETE    = 8
	IPSKB_REROUTED         = 16
)
const (
	/*
	   ZB_IDETH_IPV4  = 0x0800
	   PICO_IDETH_ARP = 0x0806
	   ZB_IDETH_IPV6  = 0x86DD

	   ZB_ARP_REQUEST   = 0x0001
	   ZB_ARP_REPLY     = 0x0002
	   ZB_ARP_HTYPE_ETH = 0x0001
	*/
	ZB_IDETH_IPV4  = 0x08
	PICO_IDETH_ARP = 0x0608
	ZB_IDETH_IPV6  = 0xDD86

	ZB_ARP_REQUEST   = 0x0100
	ZB_ARP_REPLY     = 0x0200
	ZB_ARP_HTYPE_ETH = 0x0100
)

const (
	ZB_LAYER_DATALINK  = 2 /* Ethernet only. */
	ZB_LAYER_NETWORK   = 3 /* IPv4, IPv6, ARP. Arp is there because it communicates with L2 */
	ZB_LAYER_TRANSPORT = 4 /* UDP, TCP, ICMP */
	ZB_LAYER_SOCKET    = 5 /* Socket management */
)

/* Here are some protocols. */
const (
	ZB_PROTO_IPV4  = 0
	ZB_PROTO_ICMP4 = 1
	ZB_PROTO_IGMP  = 2
	ZB_PROTO_TCP   = 6
	ZB_PROTO_UDP   = 17
	ZB_PROTO_IPV6  = 41
	ZB_PROTO_ICMP6 = 58
)
