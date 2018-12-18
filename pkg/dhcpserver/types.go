package dhcpserver

import (
        "net"
        "time"

        dhcp "github.com/krolaw/dhcp4"
)

type DHCPLease struct {
	nic    string    // Client's CHAddr
	expiry time.Time // When the lease expires
}

type DHCPHandler struct {
	ip            net.IP            // Server IP to use
	options       dhcp.Options      // Options to send to DHCP Clients
	start         net.IP            // Start of IP range to distribute
	leaseRange    int               // Number of IPs to distribute (starting from start)
	leaseDuration time.Duration     // Lease period
	leases        map[int]DHCPLease // Map to keep track of leases
}

type DHCPStaticAssignment struct {
	nic		string		// MAC address
	ip		net.IP		// assigned IP
	name		string		// hostname
}
