package main

import (
	"fmt"
	"nat-test/common/config"
	"nat-test/common/log"
	"nat-test/common/transfer"
	"nat-test/common/utils"
	"net"
	"strings"
	"time"
)

func main() {
	l := log.Log("Client Main")

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", config.Server1Addr, config.Server1Port))
	if err != nil {
		l.Fatal(err)
	}
	udp, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		l.Fatal(err)
	}
	sp := strings.Split(udp.LocalAddr().String(), ":")
	localIp := sp[0]
	localPort := sp[1]
	l.Info(fmt.Sprintf("local ip: %s:%s", localIp, localPort))
	localAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", localIp, localPort))
	if err != nil {
		l.Fatal(err)
	}

	// step1:
	buf := make([]byte, 64*1024)
	buf[0] = transfer.GetPublicIp1
	if _, err = udp.Write(buf[:1]); err != nil {
		l.Fatal(err)
	}
	_ = udp.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, err := udp.Read(buf); err != nil {
		l.Fatal(fmt.Errorf("UDP block: %v", err))
	}
	globalId := utils.Buf2Id(buf[:4])
	publicIp1 := utils.Buf2Ip(buf[4:8])
	publicPort1 := utils.Buf2Port(buf[8:10])
	l.Info(fmt.Sprintf("globalId: %d, publicIp: %s:%d", globalId, publicIp1, publicPort1))
	if publicIp1 == localIp {
		l.Info("NAT type: not nat")
		//return TODO
	}

	// step2:
	buf[0] = transfer.GetOrderPacket
	copy(buf[1:5], utils.Id2Buf(globalId))
	if _, err = udp.Write(buf[:5]); err != nil {
		l.Fatal(err)
	}
	_ = udp.Close()
	listen, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		l.Fatal(err)
	}
	_ = listen.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := listen.ReadFromUDP(buf); err == nil {
		l.Info("NAT type: full cone")
		//return TODO
	}
}
