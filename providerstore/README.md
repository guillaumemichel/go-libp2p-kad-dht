# Provider store

The provider store is responsible for keeping track of the provider records allocated to the peer, search and serve them.

## Provider Records

Define a format for the provider records so that they fit in the nested key-value datastore.

## Data store

The datastore in use (ipfs/go-datastore) doesn’t seem fit for nested key-value store. The go map implementation will be used at first, along with mutex, and we will modify it later if required. If the used memory is too large, need to check how to store some part of the data store on disk.

## LRU cache

If some data is to be stored on disk (slow down DHT, but still very OK). It would be great to implement a LRU cache mechanism to serve popular provider records.