package queryserv

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/queryservgo/packet"
	"io"
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
	pErr := q.process()
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

func (q *QueryServ) SendPacket(data packet.Packet, destination int) (err error) {
	if !q.IsConnected {
		err = fmt.Errorf("Not connected")
		return
	}
	sp := &packet.ServerPacket{}

	switch data.(type) {
	case *packet.ServerChannelMessage:
		sp.OpCode = ServerOP_ChannelMessage
		sp.Compressed = true
	case *packet.ServerWhoAll:
		sp.OpCode = ServerOP_WhoAll
	default:
		err = fmt.Errorf("Unknown packet request type:", data)
		return
	}

	sp.Buffer, err = data.Encode()
	if err != nil {
		err = fmt.Errorf("Error encoding buffer: %s", err.Error())
		return
	}

	//fmt.Printf("%#x\n", sp.Buffer)
	//fmt.Println("/\\ uncompressed")

	if sp.Compressed {
		sp.InflatedSize = len(sp.Buffer)
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write(sp.Buffer)
		w.Close()

		sp.Buffer, err = b.ReadBytes(0)
		fmt.Printf("%#x", sp.Buffer)
		fmt.Println("/\\ compressed")
		//io.ReadFull(bufio.NewReader(b), buffer)
	}

	spData, err := sp.Encode()
	if err != nil {
		fmt.Println("error making scm", err.Error())
		return
	}

	//fmt.Println(spData)

	fmt.Printf("%#x\n", spData)
	//fmt.Println(spData)
	//_, err = q.conn.Write(buffer)
	if err != nil {
		fmt.Errorf("Error writing: %s", err.Error())
		return
	}
	return
}

func (q *QueryServ) process() (err error) {

	buf := bufio.NewReader(q.conn)
	for q.IsConnected {
		data := make([]byte, 1024)
		sp := &packet.ServerPacket{}
		_, err = buf.Read(data)
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Println("Error deciphering", err.Error())
			return
		}

		err = sp.Decode(data)
		if err != nil {
			fmt.Println("Error encoding", err.Error())
			continue
		}

		err = q.recievePacket(sp)
		if err != nil {
			fmt.Printf("Error decoding server packet, Opcode: %#x: %s\n", sp.OpCode, err.Error())
			fmt.Println(data)
			continue
		}
	}
	return
}

func (q *QueryServ) recievePacket(sp *packet.ServerPacket) (err error) {
	if sp == nil {
		err = fmt.Errorf("Empty packet")
		return
	}

	//fmt.Printf("%#x\n", sp.OpCode)
	switch sp.OpCode {
	case 0x00:
		//fmt.Println("Ignoring 0 pad (this is an echo of PING is my theory, 07000000000000")
		return
	case ServerOP_Speech:
		fmt.Println("Speech", sp.Buffer)

		speech := &packet.ServerSpeech{}
		err = speech.Decode(sp.Buffer)
		if err != nil {
			err = fmt.Errorf("Error decoding speech: %s", err.Error())
			return
		}
		fmt.Println(speech)
	/*
		buf = bytes.NewBufferString(sp.Buffer)
		err = struc.Unpack(buf, speech)
		speech.From = strings.Trim(speech.From, "\x00")
		speech.To = strings.Trim(speech.To, "\x00")
		speech.Message = strings.Trim(speech.Message, "\x00")
		fmt.Printf("Status: %i, From: %s, To: %s, Message: %s, Type: %i, Misc: %v\n", speech.MinStatus, speech.From, speech.To, speech.Message, speech.Type, speech)*/
	case ServerOP_QueryServGeneric:
		fmt.Println("Found a queryServ generic, ignoring for now")
	default:
		err = fmt.Errorf("Unknown Packet Found. Size: %u, Opcode: %#x, Buffer: %s\n\n", sp.Size, sp.OpCode, sp.Buffer)
		return
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
