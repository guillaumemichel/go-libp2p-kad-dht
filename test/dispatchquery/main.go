package main

import (
	"context"
	"log"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/multiformats/go-multibase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	sd "github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	ss "github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/fakeendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	sq "github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	"github.com/libp2p/go-libp2p/core/peer"
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
	selfA := peer.ID("alpha") // peer.ID is necessary for ipfskadv1 message format
	var naddrA address.NetworkAddress = peer.AddrInfo{ID: selfA, Addrs: nil}
	rtA := simplert.NewSimpleRT(address.KadID(selfA), 2)
	endpointA := fakeendpoint.NewFakeEndpoint(clk, dispatcher)
	schedA := ss.NewSimpleScheduler(ctx, clk)
	servA := simserver.NewSimServer(rtA, endpointA)
	dispatcher.AddPeer(selfA, schedA, servA)

	// create peer B
	selfB := peer.ID("beta")
	var naddrB address.NetworkAddress = peer.AddrInfo{ID: selfB, Addrs: nil}
	rtB := simplert.NewSimpleRT(address.KadID(selfB), 2)
	endpointB := fakeendpoint.NewFakeEndpoint(clk, dispatcher)
	schedB := ss.NewSimpleScheduler(ctx, clk)
	servB := simserver.NewSimServer(rtB, endpointB)
	dispatcher.AddPeer(selfB, schedB, servB)

	// connect peer A and B
	endpointA.MaybeAddToPeerstore(naddrB, consts.PeerstoreTTL)
	rtA.AddPeer(ctx, selfB)
	endpointB.MaybeAddToPeerstore(naddrA, consts.PeerstoreTTL)
	rtB.AddPeer(ctx, selfA)

	// create find peer request
	_, bin, _ := multibase.Decode(targetBytesID)
	target := peer.ID(bin)
	req := ipfskadv1.FindPeerRequest(target)

	// dummy parameters
	resp := ipfskadv1.FindPeerResponse(target, []address.NodeID{}, endpointB)
	resChan := make(chan interface{}, 100)
	handleResp := func(ctx context.Context, s sq.QueryState, resp message.MinKadResponseMessage, c chan interface{}) sq.QueryState {

		return nil
	}
	sq.NewSimpleQuery(ctx, address.KadID(target), req, resp, 1, time.Second, endpointA, rtA, schedA, resChan, handleResp)

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
