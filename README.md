# cuteDB

A slick BTree on disk based key value store implemented in pure Go. The purpose of this project is to understand how a production ready key value store works. 

Steps to Run

```
mkdir db
go test
```

## BTree
In computer science, a B-tree is a self-balancing tree data structure that maintains sorted data and allows searches, sequential access, insertions, and deletions in logarithmic time. The B-tree is a generalization of a binary search tree in that a node can have more than two children.[1] Unlike self-balancing binary search trees, the B-tree is well suited for storage systems that read and write relatively large blocks of data, such as discs. It is commonly used in databases and file systems.


### Insertions and deletions
If the database does not change, then compiling the index is simple to do, and the index need never be changed. If there are changes, then managing the database and its index becomes more complicated.

Deleting records from a database is relatively easy. The index can stay the same, and the record can just be marked as deleted. The database remains in sorted order. If there are a large number of deletions, then searching and storage become less efficient.

Insertions can be very slow in a sorted sequential file because room for the inserted record must be made. Inserting a record before the first record requires shifting all of the records down one. Such an operation is just too expensive to be practical. One solution is to leave some spaces. Instead of densely packing all the records in a block, the block can have some free space to allow for subsequent insertions. Those spaces would be marked as if they were "deleted" records.

Both insertions and deletions are fast as long as space is available on a block. If an insertion won't fit on the block, then some free space on some nearby block must be found and the auxiliary indices adjusted. The hope is that enough space is available nearby, such that a lot of blocks do not need to be reorganized. Alternatively, some out-of-sequence disk blocks may be used.

### Advantages of B-tree usage for databases
The B-tree uses all of the ideas described above. In particular, a B-tree:

* keeps keys in sorted order for sequential traversing
* uses a hierarchical index to minimize the number of disk reads
* uses partially full blocks to speed insertions and deletions
* keeps the index balanced with a recursive algorithm
* In addition, a B-tree minimizes waste by making sure the interior nodes are at least half full. A B-tree can handle an arbitrary number of insertions and deletions.

