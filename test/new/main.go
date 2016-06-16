package main

import (
	"fmt"
	"github.com/xackery/queryservgo/packet"
	"github.com/xackery/queryservgo/queryserv"
	"time"
)

var qs *queryserv.QueryServ

func main() {
	var err error
	scm := &packet.ServerChannelMessage{
		//DeliverTo: "",
		//To:        "",
		From:      "Xuluu_[ShinTwo]",
		FromAdmin: 0,
		NoReply:   true,
		ChanNum:   5,
		GuildDBId: 0,
		Language:  0,
		Queued:    0,
		Message:   "TTEESSTTTEST",
	}

	qs = &queryserv.QueryServ{}
	go connectLoop()

	time.Sleep(1 * time.Second)
	fmt.Println("Sending ServerChannelMessage")
	err = qs.SendPacket(scm, 2)
	if err != nil {
		fmt.Println("Error sending packet", err.Error())
	}

	select {}
}

func connectLoop() {
	var err error
	for {
		err = qs.Connect()
		if err != nil {
			fmt.Println("Error with connect:", err.Error())
			return
		}
	}
}
