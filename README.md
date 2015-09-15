# YourDB
A shot at a distributed, sharded, 'folder-nested' table collections, multi-tabled key-key-value storage database.

## Distributed
YourDB is intended to be distributed on two levels. First, the tables will have the ability to be sharded and replicated across multiple machines. Secondly, the collections can be divided amongst the nodes in the cluster. Any node can respond to queries. The decision to incorporate transactions has not been made.

## Sharded
Got a large table? YourDB will spread it out across multiple nodes.

## Table Collections
Organize your tables into collections and organize your collections into parent collections. Maybe there will be permission schemes someday (or maybe not).

## Key-Key-Value (KKV) Storage
Tired to checking key prefixes and parsing out 'header' names? Let YourDB do it for you!

|  Key   |   Key    | Value  |
|--------|----------|--------|
| a1b2c3 | language |   Go   |
| a1b2c3 |  status  | Rocks! |
| 9f8e7d | language |  Java  |
|  ...   |   ...    |   ...  |
