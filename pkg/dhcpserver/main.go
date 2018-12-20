package dhcpserver

import (
	"log"
    "net"
	"fmt"
    "time"

    dhcp "github.com/krolaw/dhcp4"
    "github.com/gbraad/dhcpserver/pkg/dhcpserver/config"
)

var (
	staticAssignments	[]config.StaticAssignmentsConfigType
)

func StartServer(iface string, port int, config config.ConfigType) {
	log.Println("Setup")
	serverIP := net.IP{192, 168, 126, 1}

	staticAssignments = config.StaticAssignments

	handler := &DHCPHandler{
		ip:            serverIP,
		leaseDuration: 2 * time.Hour,
		start:         net.IP{192, 168, 126, 10},
		leaseRange:    100,
		leases:        make(map[int]DHCPLease, 10),
		options: dhcp.Options{
			dhcp.OptionSubnetMask:       []byte{255, 255, 255, 0},
			dhcp.OptionRouter:           []byte(serverIP), // Presuming Server is also your router
			dhcp.OptionDomainNameServer: []byte(serverIP), // Presuming Server is also your DNS server
		},
	}
        
	if(iface == "") {
		log.Println("Listen and serve")
		log.Fatal(listenAndServe(handler, port))
    } else {
		log.Println("Create connection")
		conn := createConnection(iface, port)
        log.Fatal(dhcp.Serve(conn, handler))
	}
}

func listenAndServe(handler dhcp.Handler, port int) error {
	listener, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	return dhcp.Serve(listener, handler)
}
