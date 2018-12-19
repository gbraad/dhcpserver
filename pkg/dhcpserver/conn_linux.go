package dhcpserver

import (
	"fmt"
        dhcp "github.com/krolaw/dhcp4"
        dhcpConn "github.com/krolaw/dhcp4/conn"
)

func createConnection(iface string, port int) dhcp.ServeConn {
        // Select interface on multi interface device - just linux for now
        conn, _  := dhcpConn.NewUDP4BoundListener(iface, fmt.Sprintf(":%d", port))
	return conn;
}
