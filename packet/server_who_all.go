package packet

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
)

type ServerWhoAll struct {
	Admin    int    `struc:"int16,little"`    //int16 admin;
	FromId   int    `struc:"uint32,little"`   //uint32 fromid;
	From     string `struc:"[64]byte,little"` //char from[64];
	Whom     string `struc:"[64]byte,little"` //char whom[64];
	Race     int    `struc:"uint16,little"`   //uint16 wrace; // FF FF = no race
	Class    int    `struc:"uint16,little"`   //uint16 wclass; // FF FF = no class
	LvlLow   int    `struc:"uint16,little"`   //uint16 lvllow; // FF FF = no numbers
	LvlHigh  int    `struc:"uint16,little"`   //uint16 lvlhigh; // FF FF = no numbers
	GmLookup int    `struc:"uint16,little"`   //uint16 gmlookup; // FF FF = not doing /who all gm
}

func (s *ServerWhoAll) Sanitize() {
	s.Admin = Clamp(s.Admin, 0, 36000)
	s.FromId = Clamp(s.FromId, 0, 36000)
	s.From = NullClean(s.From)
	s.From = StringClamp(s.From, 64)
	s.Whom = NullClean(s.Whom)
	s.Whom = StringClamp(s.Whom, 64)

	s.GmLookup = Clamp(s.GmLookup, 0, 36000)
}

func (s *ServerWhoAll) Encode() (packet []byte, err error) {
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

func (s *ServerWhoAll) Decode(packet []byte) (err error) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(packet)
	err = struc.Unpack(buf, s)
	if err != nil {
		return
	}

	s.Sanitize()
	return
}
