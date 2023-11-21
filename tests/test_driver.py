import json
import subprocess
from pathlib import Path
from typing import List, Dict, Any
from dataclasses import dataclass
from os import chdir

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

class Runner:
    def __init__(self, config: Config):
        self.config = config
        self.processes:Dict[str, subprocess.Process] = {}
        
    def gen_empty_db_json(self) -> str:
        empty_dict = {
            "TableName": "Student Courses",
            "PartitionKey": "CourseId",
            "Partitions": {}
        }
        return json.dumps(empty_dict)
    
    def initialise(self):
        # Generate config.json file
        Path("./config.json").write_text(self.config.to_json())
        print("Generated config.json")

        # Generate new db file and hh file
        for node_id in self.config.ring:
            Path(f"./dbFiles/node-{node_id}.json").write_text(self.gen_empty_db_json())
            Path(f"./hintedHandoffs/node-{node_id}.json").write_text("{}")

        # Initialise nodes
        for node_id in self.config.ring:
            print(f"Initialising node {node_id}")
            self.processes[node_id] = subprocess.Popen(["go", "run", ".", f"-port={node_id}"])
        
    def kill_node(self, node_id: str):
        #TODO: add detection for invalid node
        self.processes[node_id].kill()
        

class ConfigEncoder(json.JSONEncoder):
    def default(self, o: Config) -> Dict[str, Any]:
        return o.__dict__

if __name__ == "__main__":
    # Initialise project directory as root
    chdir("../")
    
    newConf = Config("QUORUM", 100, 10, 3, {
        "3000": Node(3000, False),
        "3001": Node(3001, False)
    })

    # Init runner
    runner = Runner(newConf)
    runner.initialise()

    input()
        
