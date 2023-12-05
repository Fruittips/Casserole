import subprocess
import threading
import os
from pathlib import Path
from time import sleep, time
from queue import Queue, Empty
from gen_config import Config, gen_system
from typing import Dict

LOG_DIR  = Path("./node_logs")
START_TIMEOUT = 30


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
            subprocess.call(['taskkill', '/F', '/T', '/PID', str(proc.pid)], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        else:
            proc.kill()
        proc.wait()
        print(f"Node {node_id}: Killed.")

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
    
    def initialise(self) -> bool:
        """ Initialises the system. Blocks until all nodes have started (i.e. the server is up for all nodes) """
        # Initialise log directory
        print(f"Runner: Log directory set at {LOG_DIR.absolute()}")
        LOG_DIR.mkdir(parents=True, exist_ok=True)
        
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
    config = gen_system(
        args.nodeCount,
        args.consistencyLevel,
        100,
        args.timeout,
        args.rf,
    )
    
    runner = None
    try:
        runner = Runner(config)
        started = runner.initialise()
        if started:
            print("System initialised.")
            input("Press ENTER or CTRL-C to exit.")
        else:
            print(f"Unable to initialise all nodes within the timeout frame of {START_TIMEOUT} seconds.")
    except KeyboardInterrupt:
        pass
    finally:
        print("Exiting.")
        if runner is not None:
            runner.exit()
