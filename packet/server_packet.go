package packet

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
)

type ServerPacket struct {
	Size    uint32  `struc:"uint32,little,sizeof=Buffer"` //uint32 size
	precode [1]byte //this is an odd padding issue
	//Opcode  uint16  `struc:"uint16,little"` //uint16 opcode
	Buffer []byte
	//Wpos         uint32 `struc:"uint32,little"` //uint32 _wpos
	//Rpos         uint32 `struc:"uint32,little"` //uint32 _rpos
	//Compressed   bool   `struc:"bool,little"`   //bool   compressed
	//InflatedSize uint32 `struc:"uint32,little"` //uint32 InflatedSize
	//Destination  uint32 `struc:"uint32,little"` //uint32 destination*/
}

func (s *ServerPacket) Sanitize() {
	/*s.DeliverTo = StringClamp(s.DeliverTo, 64)
	s.To = StringClamp(s.To, 64)
	s.From = StringClamp(s.From, 64)
	s.FromAdmin = Clamp(s.FromAdmin, 0, 256)
	s.ChanNum = Clamp(s.ChanNum, 0, 36000)
	s.GuildDBId = Clamp(s.GuildDBId, 0, 36000)
	s.Language = Clamp(s.Language, 0, 36000)
	s.Queued = Clamp(s.Queued, 0, 256)
	s.Message = StringClamp(s.Message, 511)*/
}

func (s *ServerPacket) Encode() (packet []byte, err error) {
	s.Sanitize()
	s.Size = (uint32)(len(s.Buffer))

	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 1024))
	err = struc.Pack(buf, s)
	if err != nil {
		err = fmt.Errorf("Error packing payload: %s", err.Error())
		return
	}

	packet, err = ioutil.ReadAll(buf)
	if err != nil {
		err = fmt.Errorf("erro reading buffer: %s", err.Error())
		return
	}

	return
}

func (s *ServerPacket) Decode(packet []byte) (err error) {
	return
}
