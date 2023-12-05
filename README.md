# Casserole

<p align="center">
    <img src="./frontend/logo.png"/>
</p>

As new cohort sizes grow in SUTD, the number of students bidding for modules increases. With module bidding being a period in time where the entire student population will be interacting with the database, low latency for write/read operations despite the high read and write throughput are critical for the system. Traditional SQL databases are configured in a master-slave configuration, providing high read throughput, however, write operations are still bottlenecked at the master node.

Therefore, our iteration of Cassandra, aptly dubbed "Casserole", seeks to address those limitations by prioritizing high availability, partition tolerance, and eventual consistency which are key for scenarios requiring persistent availability and the capacity to handle sudden spikes in read/write throughput.

This is important for situations that cannot afford to have their databases go offline while being able to serve high read/write throughput. Availability is guaranteed through replication and the fact that any node can be the coordinator node. With no single point of failure, there will be no downtime. Partition Tolerance is assured through consistent hashing, via the distribution of data across the cluster ring that safeguards against node failures, while maintaining data accessibility. To assure eventual consistency, Casserole by default adopts a Quorum consistency level and comes built in with read repairs, hinted handoffs and full repairs.

## Usage

To initialise the system, you need a `config.json` file to define the configuration Casserole will use, and the necessary database and hinted handoff files for each node. The dashboaard html will also be updated to reflect the total number of nodes configured.

This can be automatically generated with the following Python script:
```shell
gen_config.py [-h] [-n NODECOUNT] [-c {ONE,TWO,THREE,QUORUM,ALL}] [-t TIMEOUT] [-r RF]
```

Each node can be individually run with:
```shell
go run . -port=NODE_ID
```

### Running the Entire System

Alternatively, you can generate the system **and** run the system together with the following Python script:
```shell
run_system.py [-h] [-n NODECOUNT] [-c {ONE,TWO,THREE,QUORUM,ALL}] [-t TIMEOUT] [-r RF]
```


### Running the Webpage
We provide a local frontend webpage to interact with the system. Simply open `frontend/casseroleChef.html`.

## Testing

### Performance Tests
Performance tests are in the `performanceTest` folder. Further documentation can be found in [PERFORMANCE.md](./PERFORMANCE.md).

```
cd performanceTest
go test -ports=3000
go test -ports=3000,3001,3002,3003,3004
```

### System Tests
System tests are conducted in Python with a test driver similar to the auto-run script above. More information can be found in [tests/README.md](./tests/README.md).
```
cd tests
pip install -r requirements.txt
test_driver.py
TestWriteThenRead.py
TestMultiWriteOfDiffData.py
TestMultiWriteToDiffNodes.py
```

## Implementation Notes

### External dependencies
- Murmur3 Hashing Algorithm: (https://github.com/spaolacci/murmur3)
- HTTP Server Framework: (https://github.com/gofiber/fiber)

### Helper Functions
Located at `utils`, unless stated otherwise.

### Database Manager
The Database Manager has 3 structs. These structs describe the way our data is stored.
- **DatabaseManager**: Handles filepath, holds mutex lock, Wraps database data
- **Database**: Holds tablename, PartitionKey and stores Partitions
  - Partitions store Row data, in a map. The map key is the partition key.
- **Row**: StudentID, CreatedAt, DeletedAt, StudentName

The Database Manager handles 3 main functions. These functions are done at the lowest level, writing and reading directly from the json file our data is stored in.
- Gets row data with the partition key
- Appends rows to the database
- Creates new database managers for new nodes.

### Node Manager 
The Node Manager manages a single node. It holds information on the relevant paths to information about the node. The Node Manager struct also gives us information on the nodes
- Database Manager
- Hinted Handoff Manager
- Consistent Hash Table (CHT)
- System Configuration
- Replica nodes
- Quorum 
- Local id. 

The Node Manager has 4 main functions. 
- **Creating a new node**
- **Liveness manipulation**: IsDead, MakeDead, MakeAlive. 
- **Get other Nodes**: Conducted by port, by id, for access to keys, for access to the ports
- **Get Config**: Returns a read-only reference of the configuration of the node manager

### HTTP Client
The HTTP Client enables intra-system requests and responses between nodes. The HTTP Client has 2 functions:
- **Send Internal Read**: Sends get requests and waits for system response.
- **Send Internal Write**: Sends post request with encoded data and waits for system response

### System Config
SysConfig is used at startup to provide system-level configuration options including:
- **Consistency level**
- **Grace period**
- **Timeout**
- **Number of replicas**
- **Map of Nodes**: this is built in the nodeConfig struct, with data on port id and liveness.

It also provides a load config function to access the system configuration from a given path. These configurations can be written on startup with the config.json file.

### Hinted Handoff Manager
Hinted Handoff Manager handles failed writes to temporarily dead nodes by storing data in the node that receives the failed status of the write. It has 3 structs. 
- **HintedHandoffManager**: Containing the file path, the mutex lock, and the data
- **HintedHandoff**: High level struct storing the data. The keys to the map are the nodeIds of the dead node that should receive the data. 
- **AtomicDbMessage**: Contains the Row data to be written to the new node. It also contains the CourseId: the key of the Row to be stored.

Hinted Handoff has 3 functions:
- **Creating the new hinted handoff manager**
- **Appending data**: given the node id and the AtomicDbMessage
- **OverwriteWithMem**: Writes from the struct in memory into the hinted handoff json file.

### Consistent Hash Table
To determine which node should be used to store a given key, we use a consistent hash table (CHT), with the `murmur3` hash. In our implementation, we utilised a binary search tree for faster computation.

The key function of the CHT in relation to Casserole is the retrieval of the node ID to store a given key. Since every node has the node ID of every other node in Casserole and uses the same seed for `murmur3`, this allows a predetermined consensus on which nodes should be used to store a given key, and every node can independently identify which node to forward a request to.

*Note: More info in the [./utils/cht/README.md](./utils/cht/README.md)*

