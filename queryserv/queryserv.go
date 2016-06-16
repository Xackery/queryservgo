package queryserv

import (
	"fmt"
	"github.com/xackery/queryservgo/packet"
)

type QueryServ struct {
}

func (m *QueryServ) SendPacket(data packet.Packet, destination int) (err error) {
	sm := &packet.ServerPacket{}

	switch data.(type) {
	case *packet.ServerChannelMessage:
		//sm.Opcode = ServerOP_ChannelMessage
	default:
		err = fmt.Errorf("Unknown packet request type:", data)
		return
	}
	sm.Buffer, err = data.Encode()
	if err != nil {
		return
	}
	//fmt.Println(bBuffer)
	//sm.Buffer = string(bBuffer)

	buffer, err := sm.Encode()
	if err != nil {
		fmt.Println("error making scm", err.Error())
		return
	}
	fmt.Println(buffer)
	return
}
