package network

import (
	"log/slog"
	"net"
)

// GetLocalIP retrieves the local IP address of the machine by dialing a network connection to a remote server.
// It uses UDP protocol to dial the given server address and retrieves the local IP address from the local address of the connection.
//
// Example usage:
//
//	host := "0.0.0.0"
//	if ip, err := GetLocalIP(); err == nil {
//	    host = ip.String()
//	}
//
//	fmt.Println(host)
//
// Returns the local IP address as a net.IP and any error that occurred during the network dialing operation.
const DefaultIP = "0.0.0.0"

func GetLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("error closing network connection", "error", err.Error())
		}
	}(conn)

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}
