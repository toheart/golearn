package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	connmgr2 "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
	"time"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"io"
	"log"
	mrand "math/rand"
)

func main() {
	// LibP2P 代码使用 golog 来记录消息。他们用不同的方式记录
	//字符串 ID（即“swarm”）。我们可以控制详细程度
	// 所有记录仪具有：
	golog.SetAllLoggers(golog.LevelDebug)
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	global := flag.Bool("global", false, "use global ipfs peers for bootstrapping")
	flag.Parse()
	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}
	var bootstrapPeers []peer.AddrInfo
	var globalFlag string
	if *global {
		log.Println("using global bootstrap")
		bootstrapPeers = IPFS_PEERS
		globalFlag = " -global"
	} else {
		log.Println("using local bootstrap")
		bootstrapPeers = getLocalPeerInfo()
		globalFlag = ""
	}
	ha, err := makeRoutedHost(*listenF, *seed, bootstrapPeers, globalFlag)
	if err != nil {
		log.Fatal(err)
	}
	// 设置 stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
	if *target == "" {
		log.Println("listening for connections")
		select {} // hang forever
	}
	// 解析target获取peerId
	peerid, err := peer.Decode(*target)
	if err != nil {
		log.Fatalln(err)
	}

	// peerinfo := peer.AddrInfo{ID: peerid}
	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/echo/1.0.0")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Fatalln(err)
	}

	out, err := io.ReadAll(s)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("read reply: %q\n", out)
}

// doEcho 读出流中的数据，并且将其写回
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s\n", str)
	_, err = s.Write([]byte(str))
	return err
}

func makeRoutedHost(listenPort int, randseed int64, bootstrapPeers []peer.AddrInfo, globalFlag string) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	// 设置密钥
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	connmgr, err := connmgr2.NewConnManager(100, 150, connmgr2.WithGracePeriod(time.Minute))
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connmgr),
	}
	ctx := context.Background()

	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	// 构建数据存储（DHT 需要）。这只是一个简单的、内存中的线程安全数据存储。
	dstore := dsync.MutexWrap(ds.NewMapDatastore())
	// Make the DHT
	dht := dht.NewDHT(ctx, basicHost, dstore)
	// 创建 routed host
	routedHost := rhost.Wrap(basicHost, dht)

	// 连接所选的 ipfs 节点
	err = bootstrapConnect(ctx, routedHost, bootstrapPeers)
	if err != nil {
		return nil, err
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID()))

	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	log.Printf("Now run \"./routed-echo -l %d -d %s%s\" on a different terminal\n", listenPort+1, routedHost.ID(), globalFlag)

	return routedHost, nil
}
