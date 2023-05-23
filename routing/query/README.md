# Lookup queries

Author: [Guillaume Michel](https://github.com/guillaumemichel)

## What we need to have

For each query, we need to keep track of:
- requested KadId
- new (not queried yet) candidates, ordered by XOR distance to target
- peers that we queried already (waiting, success and timeout)