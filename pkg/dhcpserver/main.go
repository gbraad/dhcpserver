package dhcpserver

import (
	"log"
        "net"
        "time"

        dhcp "github.com/krolaw/dhcp4"
)

var (
	staticAssignments	[]DHCPStaticAssignment
	iface			string
)

const (
	port = 67
)


func StartServer() {
	log.Println("Setup")
	serverIP := net.IP{192, 168, 126, 1}

	staticAssignments = []DHCPStaticAssignment{
		DHCPStaticAssignment {
			nic:  "82:7d:df:54:21:62",
			ip:   net.IP{192, 168, 126, 11},
			name: "test1-master-0",
		},
		DHCPStaticAssignment {
			nic:  "16:91:31:2c:c2:a4",
			ip:   net.IP{192, 168, 126, 10},
			name: "test1-bootstrap",
		},
		DHCPStaticAssignment {
			nic:  "e2:14:06:fa:79:79",
			ip:   net.IP{192, 168, 126, 51},
			name: "test1-worker-0-n8dtz",
		},
	}

	handler := &DHCPHandler{
		ip:            serverIP,
		leaseDuration: 2 * time.Hour,
		start:         net.IP{192, 168, 126, 10},
		leaseRange:    50,
		leases:        make(map[int]DHCPLease, 10),
		options: dhcp.Options{
			dhcp.OptionSubnetMask:       []byte{255, 255, 255, 0},
			dhcp.OptionRouter:           []byte(serverIP), // Presuming Server is also your router
			dhcp.OptionDomainNameServer: []byte(serverIP), // Presuming Server is also your DNS server
		},
	}
        
	if(iface == "") {
		log.Println("Listen and serve")
		log.Fatal(dhcp.ListenAndServe(handler))
        } else {
		log.Println("Create connection")
		conn := createConnection(iface, port)
        	log.Fatal(dhcp.Serve(conn, handler))
	}
}
