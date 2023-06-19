# Events Management

All logic about events management including an event queue and event scheduler.

Author: [Guillaume Michel](https://github.com/guillaumemichel)

## Event queue

### Thread unsafe event queues

For simulations, it is fine to have an event queue that doesn't support multithreading.

### Priority

List of priority in order (without numbers so that it is easier to insert new events or reorder):

- ctx cancel
- IO (read from provider store)
- query cancel
- server requests (from remote peers)
- handle response to sent requests
- sending the first messages of a query (first in terms of concurrency)
- new client requests (find/provide)
- request timeout
- ...
- sending the last messages of a query
- bucket/node refresh
- reprovide