package protocol

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
)

func NewPingPong() p2p.Protocol {
	return p2p.Protocol{
		Name:    "pingpong",
		Version: 1,
		Length:  2,
		Run:     runPingPong,
	}
}

func runPingPong(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	for {
		if err := p2p.Send(rw, 0, "ping"); err != nil {
			return nil
		}
		fmt.Println("send ping")
		msg, err := rw.ReadMsg()
		if err != nil {
			return err
		}
		var recv string
		if err := msg.Decode(&recv); err != nil {
			return nil
		}
		fmt.Println("recv:", recv)
		time.Sleep(time.Second)
	}
}
