package packet

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
)

type ServerChannelMessage struct {
	DeliverTo string `struc:"[64]byte,little"` //char deliverto[64];
	To        string `struc:"[64]byte,little"` //char to[64];
	From      string `struc:"[64]byte,little"` //char from[64];
	FromAdmin int    `struc:"uint8,little"`    //uint8 fromadmin;
	NoReply   bool   `struc:"bool,little"`     //bool noreply;
	ChanNum   int    `struc:"uint16,little"`   //uint16 chan_num;
	GuildDBId int    `struc:"uint32,little"`   //uint32 guilddbid;
	Language  int    `struc:"uint16,little"`   //uint16 language;
	Queued    int    `struc:"uint8,little"`    //uint8 queued; // 0 = not queued, 1 = queued, 2 = queue full, 3 = offline
	Message   string //`struc:"[5byte,little"` //char message[0];
}

func (s *ServerChannelMessage) Sanitize() {
	s.DeliverTo = StringClamp(s.DeliverTo, 64)
	s.DeliverTo = NullClean(s.DeliverTo)
	s.To = StringClamp(s.To, 64)
	s.To = NullClean(s.To)
	s.From = StringClamp(s.From, 64)
	s.From = NullClean(s.From)
	s.FromAdmin = Clamp(s.FromAdmin, 0, 256)
	s.ChanNum = Clamp(s.ChanNum, 0, 36000)
	s.GuildDBId = Clamp(s.GuildDBId, 0, 36000)
	s.Language = Clamp(s.Language, 0, 36000)
	s.Queued = Clamp(s.Queued, 0, 256)
	s.Message = StringClamp(s.Message, 512)
	s.Message = NullClean(s.Message)
}

func (s *ServerChannelMessage) Encode() (packet []byte, err error) {
	s.Sanitize()

	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, len(s.Message)+214))
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

func (s *ServerChannelMessage) Decode(packet []byte) (err error) {
	return
}
