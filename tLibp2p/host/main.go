package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	connmgr2 "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

/**
@file:
@author: levi.Tang
@time: 2024/6/27 15:19
@description:
**/

func main() {
	run()
}
func run() {
	h, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	defer h.Close()

	fmt.Printf("Hello world, my hosts Id is %s \n ", h.ID())

	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519,
		-1,
	)

	if err != nil {
		panic(err)
	}

	var idht *dht.IpfsDHT

	connmgr, err := connmgr2.NewConnManager(100, 400, connmgr2.WithGracePeriod(time.Minute))
	if err != nil {
		panic(err)
	}
	h2, err := libp2p.New(
		// 设置连接证书
		libp2p.Identity(priv),
		// 设置监听端口
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",
			"/ip4/0.0.0.0/udp/9000/quic",
		),
		// 设置tls证书
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.Routing(func(host host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(context.Background(), h)
			fmt.Printf("found idht: %s \n", idht.Host())
			return idht, err
		}),
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}
	defer h2.Close()
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ := peer.AddrInfoFromP2pAddr(addr)
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		fmt.Printf("h2 connect: %s \n", pi.Addrs)
		h2.Connect(context.Background(), *pi)
	}
	fmt.Printf("Hello World, my second hosts ID is %s\n", h2.ID())
}
