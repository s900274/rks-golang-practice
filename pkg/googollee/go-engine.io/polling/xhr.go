package polling

import (
	"gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/googollee/go-engine.io/transport"
)

var Creater = transport.Creater{
	Name:      "polling",
	Upgrading: false,
	Server:    NewServer,
	Client:    NewClient,
}
