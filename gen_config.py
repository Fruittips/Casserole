import json
from pathlib import Path
from typing import Dict, Any
from dataclasses import dataclass

ENCODING = "utf-8"
START_NODE_ID = 3000

### CLASSES ###

@dataclass
class Node:
    port: int
    isDead: bool

@dataclass
class Config:
    consistencyLevel: str
    gracePeriod: int
    timeout: int
    rf: int
    ring: Dict[str, Node]

    def to_json(self) -> str:
        return json.dumps(self, cls=ConfigEncoder)


class ConfigEncoder(json.JSONEncoder):
    def default(self, o: Config) -> Dict[str, Any]:
        return o.__dict__


def gen_system(n: int = 5, consistencyLevel:str="QUORUM", gracePeriod=100, timeout=10, rf=3) -> Config:
    # Generate config
    nodesDict = {}
    for node_id in range(START_NODE_ID, START_NODE_ID + n):
        nodesDict[str(node_id)] = Node(node_id, False)
    config = Config(consistencyLevel, gracePeriod, timeout, rf, nodesDict)
    
    # Generate config.json file
    configFile = Path("./config.json")
    configFile.unlink(missing_ok=True)
    configFile.touch()
    configFile.write_text(config.to_json())

    # Generate new db file and hh files
    Path("./dbFiles").mkdir(exist_ok=True)
    Path("./hintedHandoffs").mkdir(exist_ok=True)
    for node_id in config.ring:
        dbFile = Path(f"./dbFiles/node-{node_id}.json")
        hhFile = Path(f"./hintedHandoffs/node-{node_id}.json")
        dbFile.unlink(missing_ok=True)
        hhFile.unlink(missing_ok=True)
        dbFile.touch()
        hhFile.touch()
        dbFile.write_text(json.dumps({
            "TableName": "Student Courses",
            "PartitionKey": "CourseId",
            "Partitions": {}
        }))
        hhFile.write_text("{}")

    print(f"Generated system with {n} nodes, RF={rf}, consistency level {consistencyLevel}.")
    return config

    
if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(
        description="Generate config file and database/hintedhandoff files for each node."
    )
    parser.add_argument("-n", "--nodeCount", type=int, default=5, help="Number of nodes in the system.")
    parser.add_argument("-c", "--consistencyLevel", type=str, default="QUORUM", help="Consistency level of the system.", choices=["ONE", "TWO", "THREE", "QUORUM", "ALL"])
    parser.add_argument("-t", "--timeout", type=int, default=10, help="Time in seconds to wait to conclude that a node is dead.")
    parser.add_argument("-r", "--rf", type=int, default=3, help="Replication factor of the system.")

    args = parser.parse_args()
    gen_system(
        args.nodeCount,
        args.consistencyLevel,
        100,
        args.timeout,
        args.rf,
    )
