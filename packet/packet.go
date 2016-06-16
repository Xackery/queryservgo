package packet

type Packet interface {
	Sanitize()
	Encode() (packet []byte, err error)
	Decode(packet []byte) (err error)
}
