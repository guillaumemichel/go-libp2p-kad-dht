# SimpleRouting

A routing module simple example.

Author: [Guillaume Michel](https://github.com/guillaumemichel)

## `FIND_PEER` RPC

Previous implementation is sync, and the peer.AddrInfo is returned once the remote peer is connected.

Can be either sync or async.

### Sync

The sync implementation should copy the behavior of the previous implementation.

An optimisation could be to try dialing the target peer once we receive its peer.AddrInfo. However, this call needs to be async (handled by scheduler). This event needs to be prioritised over other events from the query. This can be triggered by the "StopFn", the dial is done immediately, and when the answer comes back the function completes.

Another optimization could be never to query the peer directly.

### Async

An async implementation can return the peer.AddrInfo as soon as it is recieved even if the peer isn't connected.

For an async function, the caller should have the ability to cancel the query (once satisfied with results). This can be done by cancelling context, the thread waiting for results can then add an event to the event queue to cancel the query (prioritary.)