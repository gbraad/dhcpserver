package dhcpserver

import (
	"log"
        "net"
        "time"
	"math/rand"

        dhcp "github.com/krolaw/dhcp4"
)

func (handler *DHCPHandler) handleDiscover(packet dhcp.Packet, options dhcp.Options) (d dhcp.Packet) {
	log.Println("Discover")

	offeredIP := net.IPv4zero
	free, nic := -1, packet.CHAddr().String()
	log.Println("  MAC:", nic)

	// Find previous lease
	for i, v := range handler.leases {
		if v.nic == nic {
			free = i
			goto reply  // Yuck!
		}
	}

	// Static assignment
	for _, v := range staticAssignments {
		if v.nic == nic {
			offeredIP = v.ip
			goto reply  // Yuck!
		}
	}

	// Find a free lease (based on range)
	if free = handler.freeLease(); free == -1 {
		return
	}
	
reply:
	if offeredIP.Equal(net.IPv4zero) {  // no static
		offeredIP = dhcp.IPAdd(handler.start, free)
	}
	log.Println("  Reply", offeredIP)

	replyPacket := dhcp.ReplyPacket(packet, dhcp.Offer, handler.ip, offeredIP, handler.leaseDuration,
		handler.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	return replyPacket	
}

func (handler *DHCPHandler) handleRequest(packet dhcp.Packet, options dhcp.Options) (d dhcp.Packet) {
	log.Println("Request")

	if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(handler.ip) {
		return nil // Message not for this dhcp server
	}
	requestedIP := net.IP(options[dhcp.OptionRequestedIPAddress])
	log.Println("  Requested:", requestedIP)
	if requestedIP == nil {
		requestedIP = net.IP(packet.CIAddr())
	}

	if len(requestedIP) == 4 && !requestedIP.Equal(net.IPv4zero) {
		if leaseNum := dhcp.IPRange(handler.start, requestedIP) - 1; leaseNum >= 0 && leaseNum < handler.leaseRange {
			nic := packet.CHAddr().String()
			log.Println("  MAC:", nic)
			if lease, exists := handler.leases[leaseNum]; !exists || lease.nic == nic {
				handler.leases[leaseNum] = DHCPLease{ nic: nic, expiry: time.Now().Add(handler.leaseDuration) }
				log.Println("  Reply - ACK", requestedIP)
				replyPacket := dhcp.ReplyPacket(packet, dhcp.ACK, handler.ip, requestedIP, handler.leaseDuration,
					handler.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				replyPacket.AddOption(12, []byte("named"))

				return replyPacket
			}
		}
	}

	log.Println("  Reply - NAK")
	return dhcp.ReplyPacket(packet, dhcp.NAK, handler.ip, nil, 0, nil)

}

func (handler *DHCPHandler) handleReleaseDecline(packet dhcp.Packet, options dhcp.Options) (d dhcp.Packet) {
	log.Println("Release/Decline")
	nic := packet.CHAddr().String()
	log.Println("  MAC:", nic)
	for i, v := range handler.leases {
		if v.nic == nic {
			delete(handler.leases, i)
			break
		}
	}

	return nil
}

func (handler *DHCPHandler) ServeDHCP(packet dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	switch msgType {

	case dhcp.Discover:
		return handler.handleDiscover(packet, options)
	case dhcp.Request:
		return handler.handleRequest(packet, options)
	case dhcp.Release, dhcp.Decline:
		return handler.handleReleaseDecline(packet, options)
	}

	log.Println("Return nil")	
	return nil
}

func (h *DHCPHandler) freeLease() int {
	now := time.Now()
	b := rand.Intn(h.leaseRange) // Try random first
	for _, v := range [][]int{[]int{b, h.leaseRange}, []int{0, b}} {
		for i := v[0]; i < v[1]; i++ {
			if l, ok := h.leases[i]; !ok || l.expiry.Before(now) {
				return i
			}
		}
	}
	return -1
}
