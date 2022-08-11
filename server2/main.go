package main

import (
	"fmt"
	"nat-test/common/config"
	"nat-test/common/log"
	"nat-test/common/transfer"
	"nat-test/common/utils"
	"net"
)

func main() {
	l := log.Log("Server2 Main")

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("0.0.0.0:%d", config.Server2Port))
	if err != nil {
		l.Fatal(err)
	}
	udp, err := net.ListenUDP("udp4", addr)
	if err != nil {
		l.Fatal(err)
	}
	l.Info(fmt.Sprintf("start success: %s", udp.LocalAddr().String()))

	buf := make([]byte, 64*1024)
	for {
		if readLen, _, err := udp.ReadFromUDP(buf); err != nil {
			l.Error(err)
		} else if readLen > 0 {
			switch buf[0] {
			case transfer.GetOrderPacket:
				globalId := utils.Buf2Id(buf[1:5])
				clientIp := utils.Buf2Ip(buf[5:9])
				clientPort := utils.Buf2Port(buf[9:11])
				copy(buf[:4], utils.Id2Buf(globalId))
				if clientAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", clientIp, clientPort)); err == nil {
					if c, err := net.DialUDP("udp4", nil, clientAddr); err == nil {
						_, _ = c.Write(buf[:4])
						_ = c.Close()
					}
				}
			}
		}
	}
}
