# tests

A set of system tests for the overall Casserole database.

`test_driver.py` is a Python test driver that:
- Spawns an arbitrary number of nodes on different ports
- Performs tests through sending HTTP requests and expecting certain responses.

## System Tests
### `TestWriteThenRead.py`
This is a standard write and read of the database. 
1. Initialise 10 nodes
2. Push arbitrary data to a single node.
3. Read the same data from the single node, expecting matching data.

### `TestMultiWriteOfDiffData.py`
This tests with multiple writes of **different data** to the **same node**. The latest write should win, eventually.
- 'Different data' is defined as different `StudentName` values. 
- Since the writes are sent in 1 second apart, only the **latest write** should be read.

### `TestMultiWriteToDiffNodes.py`
This tests with multiple writes of **different data** to **different nodes**. Again, the latest write should win eventually.
