package virtualswitch

import (
	"fmt"
	"github.com/songgao/packets/ethernet"
	"net"
	"vswitch/pkg/common"
	"vswitch/pkg/mactable"
	"vswitch/pkg/porttable"
	"vswitch/pkg/util/log"
)

type VirtualSwitch struct {
	macTable        *mactable.MacTable
	portTable       *porttable.PortTable
	forwardingCount int
	floodingCount   int
	discardingCount int
}

func New() *VirtualSwitch {
	log.Trace("Create new VirtualSwitch")
	var s VirtualSwitch
	s.macTable = mactable.New()
	s.portTable = porttable.New()

	return &s
}

func (s *VirtualSwitch) Plugin(conn net.Conn) (int, error) {
	port, err := s.portTable.Add(conn)
	if err != nil {
		log.Info("VirtualSwitch is full")
		return -1, err
	}
	log.Debug(fmt.Sprintf("Port %d UP", port))
	return port, nil
}

func (s *VirtualSwitch) Unplug(port int) {
	log.Debug(fmt.Sprintf("Port %d DOWN", port))
	s.portTable.Delete(port)
}

func (s *VirtualSwitch) Process(frame ethernet.Frame, srcPort int) {
	s.macTable.Update(frame.Source(), srcPort)

	dstPort, err := s.macTable.Find(frame.Destination())
	if err != nil {
		log.Trace(fmt.Sprintf("VirtualSwitch Flooding: %d -> A %d Byte", srcPort, len(frame.Payload())))
		s.flood(frame, srcPort)
		return
	}

	if dstPort == srcPort {
		log.Trace(fmt.Sprintf("VirtualSwitch Discarding: %d -> X %d Byte", srcPort, len(frame.Payload())))
		s.discard(frame)
		return
	}

	log.Trace(fmt.Sprintf("VirtualSwitch Forwarding: %d -> %d %d Byte", srcPort, dstPort, len(frame.Payload())))
	s.forward(frame, dstPort)
}

func (s *VirtualSwitch) forward(frame ethernet.Frame, dstPort int) {
	s.floodingCount++
	conn, err := s.portTable.Find(dstPort)
	if err != nil {
		s.discard(frame)
		return
	}
	err = common.SendPackage(conn, frame)
	if err != nil {
		s.portTable.Delete(dstPort)
		s.discard(frame)
		return
	}
}

func (s *VirtualSwitch) flood(frame ethernet.Frame, srcPort int) {
	s.floodingCount++
	iPort := s.portTable.NextUpPort(srcPort)
	for {
		if iPort == srcPort {
			break
		}

		conn, err := s.portTable.Find(iPort)
		if err != nil {
			s.discard(frame)
			continue
		}
		err = common.SendPackage(conn, frame)
		if err != nil {
			s.portTable.Delete(iPort)
			s.discard(frame)
			continue
		}

		iPort = s.portTable.NextUpPort(iPort)
	}
}

func (s *VirtualSwitch) discard(_ ethernet.Frame) {
	s.discardingCount++
}
