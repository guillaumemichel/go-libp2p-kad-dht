# DHT Modules

Major modules of the DHT. Additional modules can be added.

```mermaid
---
title: DHT Modules
---
flowchart TD
    K["<font size=5><b>Kubo</b></font>"] --> P["<font size=5><b>Provider</b></font><br>StartProvide()<br>StopProvide()<br>ProvideList()"]
    K --> PS["<font size=5><b>Provider Store</b></font><br>Provider Records<br>Database<br>Cache (LRU)"]
    P --> PS
    P --> R
    K --> R["<font size=5><b>Routing</b></font><br>FindClosestPeers()<br>GetValue()<br>FindPeer()"]
    S["<font size=5><b>Server</b></font><br>Put()<br>Get()<br>FindClosestPeers()"] --> PS
    R -.->|over network| S
    S --> R
    S --> RT["<font size=5><b>Routing Table</b></font><br>LocalClosestPeers()"]
    R --> RT
```

I.e `fullrt` only has a different `Routing Table` implementation (allowing to bypass some `Routing`).

## DHT Modules

### Provider Store

`Provider Store` is responsible for storing and serving provider records, either pinned locally or remotely. It defines a specific `Provider Record` format. The design of the `Provider Store` database is left to the module implementer.

### Provider

`Provider` is responsible for providing and reproviding content pinned through its interface (`StartProvide` and `StopProvide`). It writes pinned content to the `Provider Store`, and makes use of `Routing` to find the appropriate peers to allocate the provider records.

### Routing Table

`Routing Table` is a database of peer identities and multiaddresses. It must expose at least a `LocalClosestPeers` function returning the closest peers in XOR distance to a provided key. `Routing Table` is responsible for deciding how many peers, and which peers are recorded.

### Routing

`Routing` is a large module, responsible for finding remote peers or values in the network. It needs `Routing Table` to keep track about the identity and multiaddresses of remote peers. `Routing` performs the iterative DHT query process.

### Server

`Server` is responsible for handling incoming requests. For `Put` and `Get` requests it must write to and read from the `Provider Store`. For `FindClosestPeers` requests, it needs the `Routing Table`'s `LocalClosestPeers` function. Actions in `Server` are usually triggered by the `Routing` modules from another peer across the network.

### Network

Use proto3 instead of custom data format.
Use r := pbio.NewDelimitedReader(s, signedIDSize), as reader and similar writer.
Use sync.Pool to avoid too much allocations.

## DHT Interface

For now, the only required interface is a `StartProvide`, `StopProvide`, `ProvideList`, `FindProvs`, `FindPeer`. For IPNS, a `GetValue` will be required.

# [Mutlitier Architecture](https://en.wikipedia.org/wiki/Multitier_architecture)

Each module should follow the three-tier architecture. The Presentation Tier contains the interfaces (high level abstractions). The Presentation Tier of each module is assumed to be the unique source of truth to what the module should be doing. The Application Tier contains the actual implementations. There can be multiple implementations for each module, they can either live in the same folder, or in another repository.

A Data Tier may be required in the future for specific modules.

## DHT Lookup state machine

The goal is for the implementation to be single threaded. The most challenging part is request concurrency, in the Routing Module. The DHT Lookup process state machine can be described as follow:

```mermaid
---
title: DHT Lookup State Machine
---
stateDiagram-v2
    direction LR
    
    Wait
    
    %% new lookup request
    DefineQueue: Define empty peer queue</br>for the new request
    FindClosestInRT: Find closest</br>peers in the RT
    AddPeersToQueue: Add peers to queue
    SendRequest: Send request to</br>first peer in queue
    MarkPeerContacted: Mark this peer as contacted

    WaitForResponses: Wait for responses
    ProviderRecordFound: Provider record found?

    ReturnResult: Return result to requester
    AllRequestsAnswered: All requests answered?

    Wait --> DefineQueue: new lookup</br>request
    DefineQueue --> FindClosestInRT
    FindClosestInRT --> AddPeersToQueue
    AddPeersToQueue --> SendRequest
    SendRequest --> MarkPeerContacted
    MarkPeerContacted --> SendRequest: #inflight messages <</br>concurrency factor
    MarkPeerContacted --> WaitForResponses: #inflight messages ==</br>concurrency factor
    WaitForResponses --> ProviderRecordFound: got a response
    WaitForResponses --> DefineQueue: new lookup</br>request
    WaitForResponses --> SendRequest: request timed out
    ProviderRecordFound --> AddPeersToQueue: No
    ProviderRecordFound --> ReturnResult: Yes
    ReturnResult --> AllRequestsAnswered
    AllRequestsAnswered --> Wait: Yes
    AllRequestsAnswered --> WaitForResponses: No

    note left of DefineQueue
        Maintain a list of ongoing requests
        Each ongoing request has a peer queue
    end note

    note left of SendRequest
        For this request, keep track of the peers
        we have contacted, and the time
    end note

    note right of MarkPeerContacted
        We don't want to contact the same peer twice
        Remove from queue, but remember we are waiting for a response
    end note

    note right of ReturnResult
        Remove request from ongoing requests list
    end note
```
