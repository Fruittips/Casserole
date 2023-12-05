from os import chdir
from time import time, sleep
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, read_req, NodeWatcher, LOG_DIR
from sys import exit
from pathlib import Path

def TestRR(runner: Runner) -> bool:
    """
    Triggers a sequence of events that causes a read repair to occur.
    """

    # Select random nodes
    write_node_id = "3000"
    outdated_node_id = "3001"

    # Test Data
    course_id = "50.069"
    student_id = "1006969"
    student_name = "TestUser1"

    print("\n\n---SIMULATING READ REPAIR.---")
    print("We write to a node, in a system of 3 nodes with RF=3 and consistency QUORUM.")
    print("We then simulate inconsistent data by clearing the database of one node.")
    print("Finally, a read request should then trigger a read repair of the affected node.")

    input("ENTER to write to node...")
    print()

    # Write to nodes
    data = get_write_data(student_id, student_name, 0, -1)
    print(f"1. Write to node {write_node_id} with data {data}.")
    resp, success = write_req(write_node_id, course_id, data)
    if not success:
        print(f"1. Error: {resp}")
        return False
    print(f"1. Write Response: {resp}")
    input("ENTER to simulate inconsistent data...")
    print()

    # Simulate inconsistent data
    print("2. Simulate inconsistent node: Shutdown process.")
    ## Kill the subprocess running the node
    runner.watchers[outdated_node_id].exit()
    ## Clear the db file and hh file
    dbFile = Path(f"./dbFiles/node-{outdated_node_id}.json")
    hhFile = Path(f"./hintedHandoffs/node-{outdated_node_id}.json")
    dbFile.unlink()
    hhFile.unlink()
    dbFile.touch()
    hhFile.touch()
    dbFile.write_text(runner.gen_empty_db_json())
    hhFile.write_text("{}")

    print("2. Simulate inconsistent node: Rebooting process...")
    ## Revive the subprocess
    runner.watchers[outdated_node_id] = NodeWatcher(outdated_node_id, LOG_DIR)
    runner.watchers[outdated_node_id].start()

    ## Wait until node started
    print(f"Runner: Initialising N{outdated_node_id}...")
    node_start_time = time()
    while not runner.watchers[outdated_node_id].startEvent.is_set():
        if runner.watchers[outdated_node_id].abortEvent.is_set():
            return False
        pass
    print(f"Runner: N{outdated_node_id} initialised after {time() - node_start_time} seconds.")
    
    input("ENTER to read (and hence trigger read repair)...")
    print()
    
    # Read from node
    print(f"3. Trigger read repair by reading from N{write_node_id}")
    resp, success = read_req(write_node_id, course_id, data["StudentId"])
    if not success:
        print(f"3. Error: {resp}")
    print(f"3. Read Response: {resp}")

    input("ENTER to exit...")
    return True

if __name__ == "__main__":
    runner = None
    try:
        # Initialise project directory as root
        chdir("../")
        config = getConfigWithNNodes(3)
        runner = Runner(config)
        started = runner.initialise()
        if not started:
            print("FAIL. Unable to initialise all nodes.")
            runner.exit()
            exit(1)

        print("---INITIALISING TEST---")
        start = time()
        success = TestRR(runner)
        if success:
            print(f"COMPLETED. Refer to Packet Capture. Time taken: {time()-start}")
        else:
            print(f"FAIL. Time taken: {time()-start}")

        print("Exiting.")
        runner.exit()
        
    except KeyboardInterrupt:
        print("Interrupted. Exiting.")
        if runner is not None:
            runner.exit()
