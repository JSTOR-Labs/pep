package netdev

import (
	"errors"
	"net"

	"github.com/rs/zerolog/log"
)

func GetMACAddress() (net.HardwareAddr, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	currentIP := conn.LocalAddr()

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	udpAddr := currentIP.(*net.UDPAddr)

	for _, i := range interfaces {
		if addrs, err := i.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					log.Info().Msgf("%s==%s", ipnet.IP.String(), udpAddr.IP.String())
					if ipnet.IP.String() == udpAddr.IP.String() {
						return i.HardwareAddr, nil
					}
				}
			}
		}
	}

	return nil, errors.New("device not found")
}
