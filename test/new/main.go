package main

import (
	"fmt"
	"github.com/xackery/queryservgo/packet"
	"github.com/xackery/queryservgo/queryserv"
)

func main() {

	scm := &packet.ServerChannelMessage{
		DeliverTo: "asaoisdjsaodiajsodijasodij",
		To:        "b",
		From:      "c",
		FromAdmin: 1,
		NoReply:   true,
		ChanNum:   3,
		GuildDBId: 4,
		Language:  5,
		Queued:    6,
		Message:   "",
	}

	qs := &queryserv.QueryServ{}
	err := qs.Connect()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	qs.SendPacket(scm, 1)
}
