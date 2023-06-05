package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/events/eventqueue/chanqueue"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	"github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting"
	sq "github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server"
	"github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multibase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	targetBytesID = "mACQIARIgp9PBu+JuU8aicuW8xT+Oa08OntMyqdLbfQtOplAHlME"
)

func lookupTest(ctx context.Context) {
	ai := serv(ctx)
	client1(ctx, ai)
}

func client1(ctx context.Context, ai peer.AddrInfo) {
	ctx, span := internal.StartSpan(ctx, "simplequerytest.client1")
	defer span.End()

	h, err := util.Libp2pHost(ctx, "9999")
	if err != nil {
		panic(err)
	}
	rt := simplert.NewSimpleRT(key.PeerKadID(h.ID()), 20)
	msgEndpoint := libp2pendpoint.NewMessageEndpoint(h)

	ai.Addrs = []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/0.0.0.0/tcp/8888")}
	h.Connect(ctx, ai)
	if !rt.AddPeer(ctx, ai.ID) {
		log.Println("failed to add peer")
	}
	ep := events.NewEventPlanner(clock.NewMock())
	eventqueue := chanqueue.NewChanQueue(ctx, 1000)
	go events.RunLoop(ctx, ep, eventqueue)

	sr, err := simplerouting.NewSimpleRouting(key.PeerKadID(h.ID()),
		msgEndpoint, rt, eventqueue, *ep)
	if err != nil {
		panic(err)
	}

	_, bin, _ := multibase.Decode(targetBytesID)
	p := peer.ID(bin)

	res, err := sr.FindPeer(ctx, p)
	if err != nil {
		panic(err)
	}
	fmt.Println("result:", res)

}

func client0(ctx context.Context, ai peer.AddrInfo) {
	ctx, span := internal.StartSpan(ctx, "simplequerytest.client0")
	defer span.End()

	h, err := util.Libp2pHost(ctx, "9999")
	if err != nil {
		panic(err)
	}
	_, bin, _ := multibase.Decode(targetBytesID)
	p := peer.ID(bin)
	marshalledPeerid, _ := p.MarshalBinary()
	msg := &pb.Message{
		Type: pb.Message_FIND_NODE,
		Key:  marshalledPeerid,
	}

	ai.Addrs = []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/0.0.0.0/tcp/8888")}

	rt := simplert.NewSimpleRT(key.PeerKadID(h.ID()), 20)
	msgEndpoint := libp2pendpoint.NewMessageEndpoint(h)

	h.Connect(ctx, ai)
	if !rt.AddPeer(ctx, ai.ID) {
		log.Println("failed to add peer")
	}

	ep := events.NewEventPlanner(clock.NewMock())
	eventqueue := chanqueue.NewChanQueue(ctx, 1000)
	go events.RunLoop(ctx, ep, eventqueue)

	resultChan := make(chan interface{}, 10)
	successFnc := func(ctx context.Context, tmp []interface{}, resp *pb.Message, resultChan chan interface{}) []interface{} {
		if len(resp.CloserPeers) > 0 {
			resultChan <- "success"
			return tmp
		}
		return tmp
	}
	sq.NewSimpleQuery(ctx, key.PeerKadID(p), msg, 1, time.Second, consts.ProtocolDHT,
		msgEndpoint, rt, eventqueue, *ep, resultChan, successFnc)

	res := <-resultChan
	switch r := res.(type) {
	case string:
		fmt.Println(r)
	}
}

func serv(ctx context.Context) peer.AddrInfo {
	ctx, span := internal.StartSpan(ctx, "simplequerytest.serv")
	defer span.End()

	h, err := util.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	em := events.NewEventsManager(ctx)
	rt := simplert.NewSimpleRT(key.PeerKadID(h.ID()), 20)
	serv := server.NewServer(ctx, h, rt, em, []protocol.ID{consts.ProtocolDHT})
	server.SetStreamHandler(serv, serv.DefaultStreamHandler, consts.ProtocolDHT)

	//p := peer.ID("12D3KooWG2qAjJvJwv4K7hrHbNVJdDzQqqwPSEezM1R3csV22yK3")
	_, bin, _ := multibase.Decode(targetBytesID)
	p := peer.ID(bin)
	h.Peerstore().AddAddr(p, multiaddr.StringCast("/ip4/1.2.3.4/tcp/5678"), peerstore.PermanentAddrTTL)
	rt.AddPeer(ctx, p)

	return h.Peerstore().PeerInfo(h.ID())
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	lookupTest(ctx)
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Kademlia-Test"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		)),
	)
	return tp, nil
}
