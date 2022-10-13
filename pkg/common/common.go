package common

import (
	"encoding/binary"
	"fmt"
	"github.com/songgao/packets/ethernet"
	"net"
	"os"
	"os/signal"
	"syscall"
	"vswitch/pkg/util/log"
)

func SendPackage(conn net.Conn, frame ethernet.Frame) error {
	length := len(frame)
	sendBuf := make([]byte, length+4)

	binary.BigEndian.PutUint16(sendBuf[0:2], uint16(0xBCBC))
	binary.BigEndian.PutUint16(sendBuf[2:4], uint16(length))

	copy(sendBuf[4:], frame)
	log.Trace(fmt.Sprintf("Send conn[%s]: [% x...] %d Byte", conn.RemoteAddr(), frame[:8], len(frame.Payload())))
	_, err := conn.Write(sendBuf)
	if err != nil {
		return err
	}

	return nil
}

func WaitExitSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		switch <-ch {
		case syscall.SIGHUP:
			log.Info(fmt.Sprintf("Recevice signal %s, ignore it", syscall.SIGINT))
		case syscall.SIGINT:
			log.Info(fmt.Sprintf("Safe exit with signal %s", syscall.SIGINT))
			return
		case syscall.SIGQUIT:
			log.Info(fmt.Sprintf("Safe exit with signal %s", syscall.SIGINT))
			return
		case syscall.SIGTERM:
			log.Info(fmt.Sprintf("Safe exit with signal %s", syscall.SIGINT))
			return
		}
	}

}
