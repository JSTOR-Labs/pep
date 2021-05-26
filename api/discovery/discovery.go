package discovery

import (
	"os"

	"github.com/grandcat/zeroconf"
)

func SetupDiscovery(port int) (*zeroconf.Server, error) {
	name, err := os.Hostname()
	if err != nil {
		name = "Unknown JSTOR Appliance"
	}
	return zeroconf.Register(name, "_jstor._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
}

func ShutdownDiscovery(server *zeroconf.Server) {
	server.Shutdown()
}
