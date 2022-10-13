package server

import (
	"fmt"
	"net"
	"os"
	"vswitch/pkg/common"
	"vswitch/pkg/config"
	"vswitch/pkg/pkgbuf"
	"vswitch/pkg/util/log"
	"vswitch/pkg/virtualswitch"
)

type Service struct {
	vSwitch *virtualswitch.VirtualSwitch

	// Accept connections from client
	listener net.Listener
}

func NewService(cfg config.ServerConf) (svr *Service, err error) {
	log.Trace("Create new Service")
	svr = &Service{
		vSwitch: virtualswitch.New(),
	}

	svr.listener, err = net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	log.Info(fmt.Sprint("Listen on: ", svr.listener.Addr()))

	return
}

func (svr *Service) Run() {
	go svr.handleListener(svr.listener)

	common.WaitExitSignal()
}

func (svr *Service) handleListener(ln net.Listener) {
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			log.Warn(fmt.Sprint(err))
		}
	}(svr.listener)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(fmt.Sprint(err))
			os.Exit(-1)
		}

		log.Info(fmt.Sprint("New connection:", conn.RemoteAddr()))

		go svr.handleConnection(conn)
	}
}

func (svr *Service) handleConnection(conn net.Conn) {
	defer func() {
		log.Info(fmt.Sprint("Close connection:", conn.RemoteAddr()))
		err := conn.Close()
		if err != nil {
			log.Warn(fmt.Sprint(err))
		}
	}()

	port, err := svr.vSwitch.Plugin(conn)
	if err != nil {
		log.Debug(fmt.Sprint(err))
		err = conn.Close()
		if err != nil {
			log.Warn(fmt.Sprint(err))
		}
	}

	defer svr.vSwitch.Unplug(port)

	log.Debug(fmt.Sprintf("Conn %s plugin port: %d", conn.RemoteAddr(), port))

	buff := pkgbuf.New(conn)
	for {
		frame, err := buff.ReadFrame()
		if err != nil {
			log.Debug(fmt.Sprint(err))
			break
		}

		log.Trace(fmt.Sprintf("Recv conn[%s]: [% x...] %d Byte", conn.RemoteAddr(), frame[:8], len(frame.Payload())))
		svr.vSwitch.Process(frame, port)
	}
}
