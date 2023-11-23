# Performance test

## Configuration
There are 2 configurations:
- Single node
  - 1 single node on port 3000 that will read and write to a json file
- Multiple nodes
  - The actual Cassandra-like implementation
  - Total of 10 nodes
  - Consistency level is Quorum
  - Replication Factor is 3

## How are tests conducted?
- There are 10 rounds of tests
- Each round will read/write the same amount of times
  - Eg. Round 1 will read 1 time and write 1 times, Round 2 will read 2 times and write 2 times, etc.
- Each round we will take the average of the time taken to perform the reads and writes
  - Average is calculated as: 
    (total time taken for all read or write requests)/(number of requests)
Note: For multiple nodes, we assume choose coordinator nodes via a round robin fashion

## Results
### Single node
Command: `go test -ports=3000`
| Round | Reads         | Writes        |
|-------|---------------|---------------|
| 1     | 1.221292ms    | 1.460625ms    |
| 2     | 438.979µs     | 795.354µs     |
| 3     | 753.264µs     | 1.077194ms    |
| 4     | 1.126552ms    | 1.146656ms    |
| 5     | 1.0023ms      | 916.358µs     |
| 6     | 1.126757ms    | 1.167916ms    |
| 7     | 1.22044ms     | 1.274893ms    |
| 8     | 1.004875ms    | 996.203µs     |
| 9     | 1.374189ms    | 1.523074ms    |
| 10    | 981.729µs     | 1.249862ms    |

### Multiple nodes (10 nodes)
Command: `go test -ports=3000,3001,3002,3003,3004,3005,3006,3007,3008,3009`
| Round | Reads       | Writes      |
|-------|-------------|-------------|
| 1     | 2.656125ms  | 2.650583ms  |
| 2     | 1.651104ms  | 1.802875ms  |
| 3     | 2.531305ms  | 2.613278ms  |
| 4     | 2.728499ms  | 2.965562ms  |
| 5     | 3.480791ms  | 3.444416ms  |
| 6     | 4.02218ms   | 4.324902ms  |
| 7     | 4.031285ms  | 4.230249ms  |
| 8     | 4.515864ms  | 5.251734ms  |
| 9     | 5.003444ms  | 5.083333ms  |
| 10    | 5.117079ms  | 5.3504ms    |
