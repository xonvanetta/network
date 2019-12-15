package network

import (
	"github.com/xonvanetta/network/client"
	"github.com/xonvanetta/network/server"
)

var (
	Client = client.New()
	Server = server.New()
)
