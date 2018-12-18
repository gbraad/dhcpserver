package dhcpserver

import (
	"log"
        "net"
        "time"
	"math/rand"

        dhcp "github.com/krolaw/dhcp4"
)

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	switch msgType {

	case dhcp.Discover:
		log.Println("Discover")

		offeredIP := net.IPv4zero
		free, nic := -1, p.CHAddr().String()
		log.Println("  MAC:", nic)

		// Find previous lease
		for i, v := range h.leases {
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
		if free = h.freeLease(); free == -1 {
			return
		}
	reply:
		if offeredIP.Equal(net.IPv4zero) {  // no static
			offeredIP = dhcp.IPAdd(h.start, free)
		}
		log.Println("  Reply", offeredIP)

		return dhcp.ReplyPacket(p, dhcp.Offer, h.ip, offeredIP, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:
		log.Println("Request")

		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
			return nil // Message not for this dhcp server
		}
		requestedIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		log.Println("  Requested:", requestedIP)
		if requestedIP == nil {
			requestedIP = net.IP(p.CIAddr())
		}

		if len(requestedIP) == 4 && !requestedIP.Equal(net.IPv4zero) {
			if leaseNum := dhcp.IPRange(h.start, requestedIP) - 1; leaseNum >= 0 && leaseNum < h.leaseRange {
				nic := p.CHAddr().String()
				log.Println("  MAC:", nic)
				if l, exists := h.leases[leaseNum]; !exists || l.nic == nic {
					h.leases[leaseNum] = DHCPLease{ nic: nic, expiry: time.Now().Add(h.leaseDuration) }
					log.Println("  Reply - ACK", requestedIP)
					return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, requestedIP, h.leaseDuration,
						h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				}
			}
		}
		log.Println("  Reply - NAK")
		return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:
		log.Println("Release/Decline")
		nic := p.CHAddr().String()
		log.Println("  MAC:", nic)
		for i, v := range h.leases {
			if v.nic == nic {
				delete(h.leases, i)
				break
			}
		}
	}
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
