package client

import (
	"fmt"
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"net"
	"os"
	"sync"
	"time"
	"vswitch/client/netcard"
	"vswitch/pkg/common"
	"vswitch/pkg/config"
	"vswitch/pkg/pkgbuf"
	"vswitch/pkg/util/log"
	"vswitch/pkg/util/util"
)

type Service struct {
	ifce    *water.Interface
	conn    net.Conn
	reconMu *sync.Mutex
	connMu  *sync.Mutex
	lAddr   *net.TCPAddr
	rAddr   *net.TCPAddr
}

func NewService(cfg config.ClientConf) (svr *Service, err error) {
	serverAddr := "172.16.1.12:8080"

	svr = &Service{}
	localAddr := netcard.SelectByUser()
	str := ""
	if util.IsIPv4(localAddr) {
		str = fmt.Sprintf("%s:0", localAddr)
	} else if util.IsIPv6(localAddr) {
		str = fmt.Sprintf("[%s]:0", localAddr)
	}
	svr.lAddr, err = net.ResolveTCPAddr("tcp", str)
	if err != nil {
		log.Error(fmt.Sprint(err))
		os.Exit(-1)
	}

	svr.rAddr, err = net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Error(fmt.Sprint(err))
		os.Exit(-1)
	}
	svr.connectServer()
	svr.ifce, err = water.New(water.Config{
		DeviceType: water.TAP,
	})
	if err != nil {
		return
	}

	log.Info(fmt.Sprint("Interface Name:", svr.ifce.Name()))

	svr.reconMu = new(sync.Mutex)
	svr.connMu = new(sync.Mutex)
	return
}

func (svr *Service) Run() {
	go svr.handleInterface(svr.ifce)
	go svr.controllerWorking()

	common.WaitExitSignal()
}

func (svr *Service) handleInterface(ifce *water.Interface) {
	var frame ethernet.Frame
	for {
		frame.Resize(1500)
		n, err := ifce.Read(frame)
		if err != nil {
			log.Error(fmt.Sprint(err))
			os.Exit(-1)
		}
		frame = frame[:n]
		log.Trace(fmt.Sprintf("Frame: %s -> %s %d Byte",
			frame.Source(), frame.Destination(), len(frame.Payload())))

		err = common.SendPackage(svr.getConn(), frame)
		if err != nil {
			log.Warn(fmt.Sprint(err))
			if svr.reconMu.TryLock() {
				svr.connectServer()
				svr.reconMu.Unlock()
			} else {
				svr.reconMu.Lock()
				svr.reconMu.Unlock()
			}

			continue
		}
	}
}

func (svr *Service) controllerWorking() {
	conn := svr.getConn()
	buff := pkgbuf.New(conn)
	for {
		frame, err := buff.ReadFrame()
		if err != nil {
			log.Warn(fmt.Sprint(err))
			if svr.reconMu.TryLock() {
				svr.connectServer()
				svr.reconMu.Unlock()
			} else {
				svr.reconMu.Lock()
				svr.reconMu.Unlock()
			}
			conn = svr.getConn()
			buff = pkgbuf.New(conn)
			continue
		}

		log.Trace(fmt.Sprintf("Recv conn[%s]: [% x...] %d Byte", conn.RemoteAddr(), frame[:8], len(frame.Payload())))
		log.Trace(fmt.Sprintf("Frame: %s <- %s %d Byte",
			frame.Destination(), frame.Source(), len(frame.Payload())))
		_, err = svr.ifce.Write(frame)
		if err != nil {
			log.Error(fmt.Sprint(err))
			os.Exit(-1)
		}
	}
}

func (svr *Service) connectServer() {
	log.Info(fmt.Sprint("Local addr:", svr.lAddr))
	log.Info(fmt.Sprint("Connect to remote addr:", svr.rAddr))
	for {
		var err error
		svr.conn, err = net.DialTCP("tcp", svr.lAddr, svr.rAddr)
		if err != nil {
			log.Warn(fmt.Sprint(err))
			log.Info("Reconnect in about 10 seconds")
			util.RandomSleep(10*time.Second, 0.9, 1.1)
			continue
		}
		break
	}
}

func (svr *Service) getConn() net.Conn {
	svr.connMu.Lock()
	defer svr.connMu.Unlock()
	return svr.conn
}
