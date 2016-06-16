package packet

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
)

type ServerSpeech struct {
	To        string `struc:"[64]byte,little"` //stringchar	to[64];
	From      string `struc:"[64]byte,little"` //char	from[64];
	GuildDBId int    `struc:"uint32,little"`   //uint32	guilddbid;
	MinStatus int    `struc:"int16,little"`    //int16	minstatus;
	Type      int    `struc:"uint32,little"`   //uint32	type;
	Message   string // `struc:"[]byte,little"`   //char	message[0];
}

func (s *ServerSpeech) Sanitize() {
	s.To = NullClean(s.To)
	s.To = StringClamp(s.To, 64)
	s.From = NullClean(s.From)
	s.From = StringClamp(s.From, 64)
	s.GuildDBId = Clamp(s.GuildDBId, 0, 36000)
	s.MinStatus = Clamp(s.MinStatus, 0, 36000)
	s.Type = Clamp(s.Type, 0, 36000)
	s.Message = NullClean(s.Message)
	s.Message = StringClamp(s.Message, 511)
}

func (s *ServerSpeech) Encode() (packet []byte, err error) {
	s.Sanitize()

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

func (s *ServerSpeech) Decode(packet []byte) (err error) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(packet)
	err = struc.Unpack(buf, s)
	if err != nil {
		return
	}
	//Since Message is variable length, we'll just take end and toss it in
	if len(packet) > 138 {
		s.Message = string(packet[138:])
	}
	s.Sanitize()
	return
}
