# Routing Table

Author: [Guillaume Michel](https://github.com/guillaumemichel)

## Challenges

2023-05-23: We want to keep track of the remote Clients that are close to us. So we want to add them in our routing table. However, we don't want to give them as _closer peers_ when answering a `FIND_NODE` request. They should remain in the RT (as long as there is space in the buckets), but not be shared.