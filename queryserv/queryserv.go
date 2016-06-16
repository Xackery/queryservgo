package queryserv

import (
	"encoding/gob"
	"fmt"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/queryservgo/packet"
	"net"
)

type QueryServ struct {
	conn        *net.TCPConn
	IsConnected bool
	config      *eqemuconfig.Config
	Addr        *net.TCPAddr
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
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", q.config.World.Tcp.Ip, q.config.World.Tcp.Port))
	if err != nil {
		err = fmt.Errorf("Error resolving address %s:%s: %s", q.config.World.Tcp.Ip, q.config.World.Tcp.Port, err.Error())
		return
	}

	fmt.Printf("Connecting to %s:%s... ", q.config.World.Tcp.Ip, q.config.World.Tcp.Port)
	q.conn, err = net.DialTCP("tcp", nil, addr)
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
	data := make([]byte, 1500)

	q.conn.Read(data) //When auth starts there's always a single echo back that can be discarded..
	pErr := q.Process()
	if err != nil {
		pErr = fmt.Errorf("Error during process: %s")
	}
	err = q.conn.Close()
	q.IsConnected = false
	if err != nil {
		err = fmt.Errorf("Error closing: %s", err.Error())
	}
	if pErr != nil || err != nil {
		if err != nil {
			err = fmt.Errorf("Error during process: %s and closing: %s", pErr.Error(), err.Error())
		}
		return
	}

	return
}

func (q *QueryServ) Process() (err error) {

	for q.IsConnected {

		sp := &packet.ServerPacket{}
		dec := gob.NewDecoder(q.conn)

		err = dec.Decode(sp)
		fmt.Printf("Got packet: %+v\n", sp)
		if err != nil {
			fmt.Printf("Error decoding server packet, Opcode: %#x: %s\n", sp.Opcode, err.Error())
			//fmt.Printf("Size: %u, Opcode: %#x, Buffer: %s\n\n", sp.Size, sp.Opcode, sp.Buffer)
			//fmt.Printf("%#X - %s\n", packet, string(packet))
			continue
		}
		err = q.recievePacket(sp)
		if err != nil {
			fmt.Printf("Error receiving packet: %s", err.Error())
			continue
		}
	}
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

func (q *QueryServ) recievePacket(sp *packet.ServerPacket) (err error) {
	if sp == nil {
		err = fmt.Errorf("Empty packet")
		return
	}
	switch sp.Opcode {
	case 0x00:
		fmt.Println("Ignoring 0 pad (this is an echo of PING is my theory, 07000000000000")
	case ServerOP_Speech:
		fmt.Println("Speech", sp.Buffer)
		/*speech := &ServerSpeech{}
		buf = bytes.NewBufferString(sp.Buffer)
		err = struc.Unpack(buf, speech)
		speech.From = strings.Trim(speech.From, "\x00")
		speech.To = strings.Trim(speech.To, "\x00")
		speech.Message = strings.Trim(speech.Message, "\x00")
		fmt.Printf("Status: %i, From: %s, To: %s, Message: %s, Type: %i, Misc: %v\n", speech.MinStatus, speech.From, speech.To, speech.Message, speech.Type, speech)*/
	default:
		fmt.Printf("Unknown Packet Found. Size: %u, Opcode: %#x, Buffer: %s\n\n", sp.Size, sp.Opcode, sp.Buffer)
	}
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
