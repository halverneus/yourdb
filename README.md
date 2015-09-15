# YourDB
A shot at a distributed, sharded, nestable table collections, multi-tabled key-key-value and object storage database.

## Distributed
YourDB is intended to be distributed on two levels. First, the tables will have the ability to be sharded and replicated across multiple machines. Secondly, the collections can be divided amongst the nodes in the cluster. Any node can respond to queries. The decision to incorporate transactions and regional replication has not been made, yet.

## Sharded
Got a large table? YourDB will spread it out across multiple nodes. Got many tables? Why not spread those out, too? Got milk? Can't help you.

## Table Collections
Organize your tables into collections and organize your collections into parent collections. Maybe there will be optional permission schemes someday (or maybe not).

## Key-Key-Value (KKV) Storage
Tired to checking key prefixes and parsing out 'column header' names? Let YourDB do it for you!

|  Key   |   Key    | Value  |
|--------|----------|--------|
| a1b2c3 | language |   Go   |
| a1b2c3 |  status  | Rocks! |
| 9f8e7d | language |  Java  |
|  ...   |   ...    |   ...  |

## Use What You Need (UWYN). You win!
Need object storage only? It is here. Want it distributed? Add one of the distribution links or make your own. Need to embed a single KKV table, or just a few? YourDB is designed from the ground up to truly be YOUR database, whatever your scale!

## Too Good To Be True?
Probably... but let's try, anyway!
