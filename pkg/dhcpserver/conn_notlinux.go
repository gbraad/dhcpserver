// +build !linux

package dhcpserver

import (
	"fmt"
        dhcp "github.com/krolaw/dhcp4"
        dhcpConn "github.com/krolaw/dhcp4/conn"
)

func createConnection(iface string, port int) dhcp.ServeConn {
        // Work around for other OSes
        conn, _ := dhcpConn.NewUDP4FilterListener(iface, fmt.Sprintf(":%d", port))

	return conn;
}
