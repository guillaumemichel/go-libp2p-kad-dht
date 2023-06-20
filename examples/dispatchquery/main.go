package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multibase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	sd "github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	ss "github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/addrinfo"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/fakeendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	sq "github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server/simipfsserver"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	peerstoreTTL = 10 * time.Minute
)

var (
	targetBytesID = "mACQIARIgp9PBu+JuU8aicuW8xT+Oa08OntMyqdLbfQtOplAHlME"
)

func queryTest(ctx context.Context) {
	ctx, span := util.StartSpan(ctx, "queryTest")
	defer span.End()

	clk := clock.NewMock()

	dispatcher := sd.NewSimpleDispatcher(clk)

	// create peer A
	selfA := peerid.PeerID{ID: peer.ID("alpha")} // peer.ID is necessary for ipfskadv1 message format
	addrA := multiaddr.StringCast("/ip4/1.1.1.1/tcp/4001/")
	var naddrA address.NetworkAddress = addrinfo.AddrInfo{
		AddrInfo: peer.AddrInfo{ID: selfA.ID, Addrs: []multiaddr.Multiaddr{addrA}}}
	rtA := simplert.NewSimpleRT(selfA.Key(), 2)
	endpointA := fakeendpoint.NewFakeEndpoint(selfA, dispatcher)
	schedA := ss.NewSimpleScheduler(ctx, clk)
	servA := simipfsserver.NewSimServer(rtA, endpointA)
	dispatcher.AddPeer(selfA, schedA, servA)

	// create peer B
	selfB := peerid.PeerID{ID: peer.ID("beta")}
	addrB := multiaddr.StringCast("/ip4/2.2.2.2/tcp/4001/")
	var naddrB address.NetworkAddress = addrinfo.AddrInfo{
		AddrInfo: peer.AddrInfo{ID: selfB.ID, Addrs: []multiaddr.Multiaddr{addrB}}}
	rtB := simplert.NewSimpleRT(selfB.Key(), 2)
	endpointB := fakeendpoint.NewFakeEndpoint(selfB, dispatcher)
	schedB := ss.NewSimpleScheduler(ctx, clk)
	servB := simipfsserver.NewSimServer(rtB, endpointB)
	dispatcher.AddPeer(selfB, schedB, servB)

	// create peer C
	selfC := peerid.PeerID{ID: peer.ID("gamma")}
	addrC := multiaddr.StringCast("/ip4/3.3.3.3/tcp/4001/")
	var naddrC address.NetworkAddress = addrinfo.AddrInfo{
		AddrInfo: peer.AddrInfo{ID: selfC.ID, Addrs: []multiaddr.Multiaddr{addrC}}}
	rtC := simplert.NewSimpleRT(selfC.Key(), 2)
	endpointC := fakeendpoint.NewFakeEndpoint(selfC, dispatcher)
	schedC := ss.NewSimpleScheduler(ctx, clk)
	servC := simipfsserver.NewSimServer(rtC, endpointC)
	dispatcher.AddPeer(selfC, schedC, servC)

	// connect peer A and B
	endpointA.MaybeAddToPeerstore(ctx, naddrB, peerstoreTTL)
	rtA.AddPeer(ctx, selfB)
	endpointB.MaybeAddToPeerstore(ctx, naddrA, peerstoreTTL)
	rtB.AddPeer(ctx, selfA)

	// connect peer B and C
	endpointB.MaybeAddToPeerstore(ctx, naddrC, peerstoreTTL)
	rtB.AddPeer(ctx, selfC)
	endpointC.MaybeAddToPeerstore(ctx, naddrB, peerstoreTTL)
	rtC.AddPeer(ctx, selfB)

	// create find peer request
	_, bin, _ := multibase.Decode(targetBytesID)
	target := peerid.PeerID{ID: peer.ID(bin)}
	req := ipfskadv1.FindPeerRequest(target)
	resp := &ipfskadv1.Message{}

	// dummy parameters
	handleResp := func(ctx context.Context, s sq.QueryState, _ address.NodeID, resp message.MinKadResponseMessage) sq.QueryState {
		fmt.Println(resp.CloserNodes())
		return nil
	}
	sq.NewSimpleQuery(ctx, target.Key(), req, resp, 1, time.Second, endpointA,
		rtA, schedA, handleResp)

	// run simulation
	dispatcher.DispatchLoop(ctx)
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

	queryTest(ctx)
}
