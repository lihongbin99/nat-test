package main

import (
	"fmt"
	"nat-test/common/config"
	"nat-test/common/log"
	"nat-test/common/transfer"
	"nat-test/common/utils"
	"net"
)

var (
	globalID int32 = 0
)

func main() {
	l := log.Log("Server1 Main")

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("0.0.0.0:%d", config.Server1Port))
	if err != nil {
		l.Fatal(err)
	}
	addr2, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", config.Server2Addr, config.Server2Port))
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
		if readLen, clientAddr, err := udp.ReadFromUDP(buf); err != nil {
			l.Error(err)
		} else if readLen > 0 {
			switch buf[0] {
			case transfer.GetPublicIp1:
				globalID++
				copy(buf[:4], utils.Id2Buf(globalID))
				copy(buf[4:8], utils.Ip2Buf(clientAddr.IP.String()))
				copy(buf[8:10], utils.Port2Buf(uint16(clientAddr.Port)))
				_, _ = udp.WriteToUDP(buf[:10], clientAddr)
			case transfer.GetOrderPacket:
				if s2, err := net.DialUDP("udp4", nil, addr2); err == nil {
					copy(buf[5:9], utils.Ip2Buf(clientAddr.IP.String()))
					copy(buf[9:11], utils.Port2Buf(uint16(clientAddr.Port)))
					_, _ = s2.Write(buf[:11])
					_ = s2.Close()
				}
			}
		}
	}
}
