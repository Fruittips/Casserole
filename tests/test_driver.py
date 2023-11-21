import json
import subprocess
import threading
from pathlib import Path
from typing import List, Dict, Any
from dataclasses import dataclass
from time import sleep

ENCODING = "utf-8"

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

class NodeWatcher:
    """ Watches the output of a process. """
    def __init__(self, node_id: str):
        self.node_id = node_id
        self.startEvent = None
        self.exitEvent = None
        self._t = None

    def start(self):
        self.startEvent = threading.Event()
        self.exitEvent = threading.Event()
        self._t = threading.Thread(target=self.thread_fn, args=(self.node_id, self.startEvent, self.exitEvent))
        self._t.start()

    @staticmethod
    def thread_fn(node_id, sEv, eEv):
        proc = subprocess.Popen(["go", "run", ".", f"-port={node_id}"],
                                bufsize=1,
                                text=True,
                                stdin=None,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT
                                )

        checkIntv = 5
        while not eEv.is_set():
            timer = threading.Timer(checkIntv, eEv.set)
            timer.start()
            try:
                out = proc.stdout.readline()
                if "initialised" in out:
                    sEv.set()
            except Exception:
                break
            timer.cancel()
            
        # End of loop
        print(f"Node {node_id}: Killed.")
        proc.kill()

    def exit(self):
        self.exitEvent.set()
        if self._t is not None:
            self._t.join()
        

class Runner:
    """ Runs the System. """
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
        """ Initialises the system. Blocks until all nodes have started (i.e. the server is up for all nodes) """
        # Generate config.json file
        Path("./config.json").write_text(self.config.to_json())
        print("Generated config.json")

        # Generate new db file and hh file
        for node_id in self.config.ring:
            Path(f"./dbFiles/node-{node_id}.json").write_text(self.gen_empty_db_json())
            Path(f"./hintedHandoffs/node-{node_id}.json").write_text("{}")

        # Initialise nodes
        for node_id in self.config.ring:
            self.watchers[node_id] = NodeWatcher(node_id)
            self.watchers[node_id].start()

        # Wait until nodes started
        for node_id in self.watchers:
            while not self.watchers[node_id].startEvent.is_set():
                pass
            print(f"Node {node_id}: Initialised.")
            
    def exit(self):
        for node_id in self.watchers:
            self.watchers[node_id].exit()

class ConfigEncoder(json.JSONEncoder):
    def default(self, o: Config) -> Dict[str, Any]:
        return o.__dict__



### TEST HELPER FUNCTIONS ###

START_NODE_ID = 3000
URL = "http://127.0.0.1"

def get_read_url(node_id: str, course_id: str, student_id: str) -> str:
    return f"{URL}:{node_id}/read/course/{course_id}/student/{student_id}"

def get_write_url(node_id: str, course_id: str) -> str:
    return f"{URL}:{node_id}/write/course/{course_id}"

def get_write_data(student_id: str, student_name: str, created_at: int, deleted_at: int) -> Dict[str, Any]:
    if deleted_at == -1:
        return {
            "StudentId": student_id,
            "StudentName": student_name,
            "CreatedAt": created_at,
            "DeletedAt": None,
        }
        
    return {
        "StudentId": student_id,
        "StudentName": student_name,
        "CreatedAt": created_at,
        "DeletedAt": deleted_at,
    }

def getConfigWithNNodes(n: int, consistencyLevel="QUORUM", gracePeriod=10, timeout=10, rf=3) -> Config:
    nodesDict = {}
    for node_id in range(START_NODE_ID, START_NODE_ID + n):
        nodesDict[str(node_id)] = Node(node_id, False)
    return Config(consistencyLevel, gracePeriod, timeout, rf, nodesDict)
