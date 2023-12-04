import json
import subprocess
import threading
import requests
import os
from pathlib import Path
from typing import Dict, Any
from dataclasses import dataclass
from time import sleep, time
from queue import Queue, Empty

ENCODING = "utf-8"
LOG_DIR  = Path("./tests/logs")
START_TIMEOUT = 30 # timeout if node fails to start after this number of seconds
# note: based on testing, ~150 should be safe for 100 nodes, ~30 is safe for 10 nodes

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
    def __init__(self, node_id: str, log_dir: Path):
        if not log_dir.is_dir():
            raise RuntimeError(f"Cannot create {node_id} log file: Log directory {log_dir} does not exist. ")
        
        self.node_id = node_id
        self.startEvent = None
        self.exitEvent = None
        self._t = None

        # Init log file
        log_file = log_dir.joinpath(f"log-{node_id}.txt")
        log_file.unlink(missing_ok=True)

        self.log_fp = log_file.open('wb')

    def start(self):
        self.startEvent = threading.Event()
        self.abortEvent = threading.Event()
        self.exitEvent = threading.Event()
        self._t = threading.Thread(target=self.thread_fn, args=(self.node_id, self.startEvent, self.abortEvent, self.exitEvent, self.log_fp))
        self._t.start()

    @staticmethod
    def thread_fn(node_id, sEv, aEv, eEv, logFp):
        str_to_detect = b"initialised"
        output_queue = Queue()
        proc = subprocess.Popen(["go", "run", ".", f"-port={node_id}"],
                                stdin=None,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.STDOUT
                                )
        proc_out_thread = threading.Thread(target=NodeWatcher.output_queue_push, args=(proc.stdout, output_queue))
        proc_out_thread.daemon = True
        proc_out_thread.start()

        checkIntv = 1
        start_attempts = START_TIMEOUT
        while not eEv.is_set():
            if not sEv.is_set():
                # Read line without block
                line_b = b""
                try:
                    line_b = output_queue.get_nowait().strip()
                except Empty:
                    pass

                if len(line_b) > 0:
                    logFp.write(line_b + b"\n")
                    if str_to_detect in line_b:
                        sEv.set()
                        continue

                start_attempts -= 1
                if start_attempts == 0:
                    print(f"Node {node_id}: Failed to start. Force quitting.")
                    aEv.set()
                    break
                sleep(checkIntv)
                continue
    
            line_b = b""
            try:
                line_b = output_queue.get_nowait().strip()
            except Empty:
                pass

            if len(line_b) > 0:
                logFp.write(line_b + b"\n")
            
            sleep(checkIntv)

        # End of loop
        if os.name == "nt":
            # Windows-specific kill logic because somehow subprocess.kill() only kills one of the processes
            subprocess.call(['taskkill', '/F', '/T', '/PID', str(proc.pid)], stdout=None)
        else:
            proc.kill()
        returncode = proc.wait()
        print(f"Node {node_id}: Killed with {returncode}.")

        # Close log file
        logFp.close()

    @staticmethod
    def output_queue_push(out, queue: Queue):
        for line in iter(out.readline, ""):
            queue.put(line)
        out.close()

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
    
    def initialise(self) -> bool:
        """ Initialises the system. Blocks until all nodes have started (i.e. the server is up for all nodes) """
        # Initialise log directory
        print(f"Runner: Log directory set at {LOG_DIR.absolute()}")
        LOG_DIR.mkdir(parents=True, exist_ok=True)
        
        # Generate config.json file
        configFile = Path("./config.json")
        configFile.unlink(missing_ok=True)
        configFile.touch()
        configFile.write_text(self.config.to_json())
        print("Runner: Generated config.json")

        # Generate new db file and hh file
        Path("./dbFiles").mkdir(exist_ok=True)
        Path("./hintedHandoffs").mkdir(exist_ok=True)
        for node_id in self.config.ring:
            dbFile = Path(f"./dbFiles/node-{node_id}.json")
            hhFile = Path(f"./hintedHandoffs/node-{node_id}.json")
            dbFile.unlink(missing_ok=True)
            hhFile.unlink(missing_ok=True)
            dbFile.touch()
            hhFile.touch()
            dbFile.write_text(self.gen_empty_db_json())
            hhFile.write_text("{}")
        return self.start()

    def start(self) -> bool:
        # Initialise nodes
        for node_id in self.config.ring:
            self.watchers[node_id] = NodeWatcher(node_id, LOG_DIR)
            self.watchers[node_id].start()

        # Wait until nodes started
        successful_start = True
        print(f"Runner: Initialising {len(self.config.ring)} nodes...")
        node_start_time = time()
        for node_id in self.watchers:
            while not self.watchers[node_id].startEvent.is_set():
                if self.watchers[node_id].abortEvent.is_set():
                    successful_start = False
                    break
                pass
            if not successful_start:
                print(f"Runner: Node {node_id} failed to initialise after {time() - node_start_time} seconds.")
            else:
                print(f"Runner: Node {node_id} initialised after {time() - node_start_time} seconds.")
        return successful_start
            
    def exit(self):
        print("Runner: Exit")
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

def get_kill_url(node_id: str) -> str:
    return f"{URL}:{node_id}/internal/kill"

def get_revive_url(node_id: str) -> str:
    return f"{URL}:{node_id}/internal/revive"

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

def read_req(svr_id: str, course_id: str, student_id: str) -> (str, bool):
    try:
        resp = requests.get(get_read_url(svr_id, course_id, student_id))
        if resp.status_code != 200:
            return f"HTTP Error {resp.status_code}, resp.text", False
        return resp.text, True
    except requests.exceptions.Timeout:
        return "Timeout", False
    except requests.exceptions.TooManyRedirects:
        return "Too Many Redirects", False
    except requests.exceptions.RequestException as e:
        return f"Catastrophic Error {e}", False

def write_req(svr_id: str, course_id: str, data: Dict[str, Any]) -> (str, bool):
    try:
        resp = requests.post(get_write_url(svr_id, course_id), data=data)
        if resp.status_code != 200:
            return f"HTTP Error {resp.status_code}, resp.text", False
        return resp.text, True
    except requests.exceptions.Timeout:
        return "Timeout", False
    except requests.exceptions.TooManyRedirects:
        return "Too Many Redirects", False
    except requests.exceptions.RequestException as e:
        return f"Catastrophic Error {e}", False

def kill_req(svr_id: str) -> bool:
    try:
        resp = requests.get(get_kill_url(svr_id))
        if resp.status_code != 200:
            return False
        return True
    except requests.exceptions.Timeout:
        return False
    except requests.exceptions.TooManyRedirects:
        return False
    except requests.exceptions.RequestException:
        return False

def revive_req(svr_id: str) -> bool:
    try:
        resp = requests.get(get_revive_url(svr_id))
        if resp.status_code != 200:
            return False
        return True
    except requests.exceptions.Timeout:
        return False
    except requests.exceptions.TooManyRedirects:
        return False
    except requests.exceptions.RequestException:
        return False
