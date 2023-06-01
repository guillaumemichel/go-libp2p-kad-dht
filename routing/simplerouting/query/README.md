# Lookup queries

Author: [Guillaume Michel](https://github.com/guillaumemichel)

## What we need to have

For each query, we need to keep track of:
- requested KadId
- new (not queried yet) candidates, ordered by XOR distance to target
- peers that we queried already (waiting, success and timeout)

## Lookup process

The lookup implementation should be the same for all RPCs. The termination condition & messages can be different.

- need to keep track of ongoing queries in a list (at the scheduler level?) so that they aren't dereferenced. Make sure they are removed once they are done/cancelled.

### `FIND_PEER`

- compute the KadID associated with the target peerid.
- initialize an empty list of peers to be queried.
    - linked list with 1 extra pointer indicating the next peer to be queried
    - each peer should have a status (queued, waiting for response, queried, timeout)
    - each peer should have its distance (XOR distance between peer's KadID and target KadID)
    - peers should be ordered according to their XOR distance to target KadID
- get the 20 closest peers from the routing table, add them to the list marked as "queued"
- add `concurrency` "send request" events to the event queue, requesting to send a request for this query. these events must contain a reference to the query
- return

- once a message to be sent is picked up by the worker
- if query has been cancelled (ctx cancel, query done) -> return
- determine appropriate remote peer (closest "queued" peer)
- send message to peer (new go routine)
- mark the peer as "waiting"
- add a request timeout to the scheduler (reference timeout event on the query)
- return

- once a response comes back (from message go routine)
- add it to the event queue (with query id / pointer)
- return (end go routine)

- when a response is picked up by the worker
- remove timeout from scheduler
- if response fulfills the success condition, report answer (and close query) -> this step must be specific for each kind of request -- RPC SPECIFIC FOLLOW-UP?
- change the peer status from "waiting" to "queried"
- try to add peer to rt and peerstore
- add (all) closer peers to the peer list (we don't want duplicates though), status "queued"
- add "send request" to the event queue for this query. ALTERNATIVELY: select the next closest "queued" peer and send it a request
- return

- when a timeout occurs
/!\ maybe an answer is arrived, but the timeout was juged prioritary (?)
- if query has been cancelled (ctx cancel, query done) -> return
- mark peer as "timeout"
- (remove peer from rt)
- add "send request" to the event queue for this query. ALTERNATIVELY: select the next closest "queued" peer and send it a request

- closing the query
- mark the query as cancelled (either with bool variable, or set peer list to nil)
- remove query from the list of queries
- return

- on failure (all peers have been queried/timeout), but no success
- if the list only contains the initial 20 peers from the routing table, increase the number of closest peers in the routing table, until one remote peer responds with a closer peer or until the whole routing table is in the query and is unresponsive.
- if remote peers have answered the query, return an error: "content not found"
- in specific conditions where a quorum is required, return an answer with a warning that the quorum wasn't reached

customs functions as parameters
- StopFn(Query): decide when to stop
- SaveUsefulState(Query): save useful responses. Can also add prioritary event to try to connect to peer
- FollowUpFn(Query): for `FIND_PEER`, once an answer is received, if it contains the peer.AddrInfo of the target peer

- following the follow up (if any)

## To explore

- `queryPeerFilter`