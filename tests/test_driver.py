import json
import subprocess
import threading
from pathlib import Path
from typing import List, Dict, Any
from dataclasses import dataclass
from os import chdir
from time import sleep

ENCODING = "utf-8"

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

class NodeWatcher:
    def __init__(self, node_id: str):
        self.node_id = node_id
        self.startEvent = None
        self.exitEvent = None

    def start(self):
        self.startEvent = threading.Event()
        self.exitEvent = threading.Event()
        t = threading.Thread(target=self.thread_fn, args=(self.node_id, self.startEvent, self.exitEvent))
        t.start()

    @staticmethod
    def thread_fn(node_id, sEv, eEv):
        proc = subprocess.Popen(["go", "run", ".", f"-port={node_id}"],
                                stdin=None,
                                bufsize=1,
                                text=True,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT
                                )

        # Watch output every x seconds
        check_interval = 0.5
        attempts = 100

        print("starting")

        out = ""

        while not eEv.is_set():
            if not sEv.is_set():
                if attempts == 0:
                    eEv.set()

                out += proc.stdout.read(encoding=ENCODING)

                if "initialised" in out:
                    print("started")
                    sEv.set()
                else:
                    attempts-=1

                continue
            
            out += proc.stdout.read(encoding=ENCODING)
                
            sleep(check_interval)

        # End of loop
        proc.kill()

    def exit(self):
        self.exitEvent.set()
        

class Runner:
    def __init__(self, config: Config):
        self.config = config
        self.watchers:Dict[str, NodeWatcher] = {}
        
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
            self.watchers[node_id] = NodeWatcher(node_id)
            self.watchers[node_id].start()

        # Wait until nodes started
        for node_id in self.watchers:
            while not self.watchers[node_id].startEvent.is_set():
                pass
            print(f"Node {node_id}: Started.")
            
    def exit(self):
        for node_id in self.watchers:
            self.watchers[node_id].exit()
    
        

class ConfigEncoder(json.JSONEncoder):
    def default(self, o: Config) -> Dict[str, Any]:
        return o.__dict__

if __name__ == "__main__":
    try:
        # Initialise project directory as root
        chdir("../")
    
        newConf = Config("QUORUM", 100, 10, 3, {
            "3000": Node(3000, False),
            "3001": Node(3001, False)
        })

        # Init runner
        runner = Runner(newConf)
        runner.initialise()
    except KeyboardInterrupt:
        print("Interrupted.")
        runner.exit()

    
        
