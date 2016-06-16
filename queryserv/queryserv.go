package queryserv

import (
	"fmt"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/queryservgo/packet"
	"net"
)

type QueryServ struct {
	conn        net.Conn
	IsConnected bool
	config      *eqemuconfig.Config
}

func (q *QueryServ) LoadConfig() (err error) {
	if q.config != nil {
		return
	}
	q.config, err = eqemuconfig.GetConfig()
	if err != nil {
		err = fmt.Errorf("Error loading config: %s", err.Error())
		return
	}
	return
}

func (q *QueryServ) Connect() (err error) {
	q.LoadConfig()

	q.IsConnected = false
	fmt.Printf("Connecting to %s:%s... ", q.config.World.Tcp.Ip, q.config.World.Tcp.Port)
	q.conn, err = net.Dial("tcp", q.config.World.Tcp.Ip+":"+q.config.World.Tcp.Port)
	if err != nil {
		fmt.Errorf("Error connecting: %s", err.Error())
		return
	}

	_, err = q.conn.Write(makeModePacket())
	if err != nil {
		fmt.Errorf("Error writing mode packet: %s", err.Error())
		return
	}
	_, err = q.conn.Write(makeAuthPacket())
	if err != nil {
		fmt.Errorf("Error writing auth packet: %s", err.Error())
		return
	}
	fmt.Println("Success!")
	q.IsConnected = true
	return
}

func (q *QueryServ) SendPacket(data packet.Packet, destination int) (err error) {
	sp := &packet.ServerPacket{}

	switch data.(type) {
	case *packet.ServerChannelMessage:
		//sm.Opcode = ServerOP_ChannelMessage
	default:
		err = fmt.Errorf("Unknown packet request type:", data)
		return
	}
	sp.Buffer, err = data.Encode()
	if err != nil {
		return
	}
	//fmt.Println(bBuffer)
	//sm.Buffer = string(bBuffer)

	buffer, err := sp.Encode()
	if err != nil {
		fmt.Println("error making scm", err.Error())
		return
	}
	fmt.Println(buffer)
	return
}

func makeModePacket() (modePacket []byte) {
	modePacket = make([]byte, 18)
	modePacket[0] = 0
	out := []byte("**PACKETMODEQS**")
	for i, b := range out {
		modePacket[i+1] = b
	}
	modePacket[17] = 0x0d
	return
}

func makeAuthPacket() (authPacket []byte) {
	authPacket = make([]byte, 23)
	authPacket[0] = 0x17
	authPacket[1] = 0x0
	authPacket[2] = 0x0
	authPacket[3] = 0x0
	authPacket[4] = 0x0  //this is ServerOP_ZAuth
	authPacket[5] = 0x25 //this is ServerOP_ZAuth
	authPacket[6] = 0x0
	authPacket[7] = 0x6c
	authPacket[8] = 0xe6
	authPacket[9] = 0x0a
	authPacket[10] = 0xcd
	authPacket[11] = 0xde
	authPacket[12] = 0x9f
	authPacket[13] = 0x9d
	authPacket[14] = 0x47
	authPacket[15] = 0x3c
	authPacket[16] = 0x4c
	authPacket[17] = 0x06
	authPacket[18] = 0x68
	authPacket[19] = 0xe5
	authPacket[20] = 0x05
	authPacket[21] = 0x54
	authPacket[22] = 0x09
	return
}
