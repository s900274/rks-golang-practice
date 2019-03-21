package websocket

import (
	"gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/googollee/go-engine.io/transport"
)

var Creater = transport.Creater{
	Name:      "websocket",
	Upgrading: true,
	Server:    NewServer,
	Client:    NewClient,
}
