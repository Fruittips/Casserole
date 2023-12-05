import json
from random import sample, choice
from os import chdir
from time import time, sleep
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, read_req
from sys import exit

def TestConsistencyAll(runner: Runner) -> bool:
    """
    Tests with multiple writes of different data to DIFFERENT NODES. The latest write SHOULD win (eventually).

    We do this by changing the name, since the things that are overwritten are:
    - Name
    - Timestamp Written
    - Timestamp Deleted
    """

    # Test Config
    sample_count = 5
    sample_reader_count = 6 # Number of nodes to read from
    write_interval = 1 # Write every {write_interval} seconds
    read_delay = 5 # Read the value {read_delay} after the writes
    read_attempts = 2  # Attempts to read until test is considered as failure
    read_interval = 1  # Read every {read_interval} seconds

    # Select random nodes
    write_node_ids = sample((list(runner.config.ring.keys())), sample_count) # nodes to write to
    read_node_ids = sample((list(runner.config.ring.keys())), sample_reader_count)

    # Test Data
    course_id = "50.069"
    student_id = "1006969"
    student_name_base = "TestUser"
    expected_name = f"{student_name_base}{sample_count}"

    # Write to nodes
    step = 1
    for node_id in write_node_ids:
        data = get_write_data(student_id, f"{student_name_base}{step}", 0, -1)
        print(f"{step}. Write to node {node_id} with data {data}.")

        resp, success = write_req(node_id, course_id, data)
        if not success:
            print(f"{step}. Error: {resp}")
            return False
        print(f"{step}. Write Response: {resp}")
        step += 1
        sleep(write_interval)


    # Read from node
    print(f"{step}. Sleeping for {read_delay} seconds.")
    sleep(read_delay)
    for node_id in read_node_ids:
        success = False
        for i in range(read_attempts):
            print(f"{step}. Read from node {node_id}")
            resp, success = read_req(node_id, course_id, student_id)
            print(f"{step}. Read Response: {resp}")

            # Check response
            read_dat = json.loads(resp)
            success = read_dat["StudentName"] == expected_name
            if success:
                break
            print(f"{step}. Error: Record read was NOT the latest. Retrying after {read_interval} seconds to account for read-repair...")
            step += 1
            sleep(read_interval)

            if not success:
                print(f"{step}. Error: Record read was NOT the latest after two tries.")
                return False
        step += 1
    return True

if __name__ == "__main__":
    runner = None
    try:
        # Initialise project directory as root
        chdir("../")
        config = getConfigWithNNodes(10, consistencyLevel="ALL", rf=10)
        runner = Runner(config)
        started = runner.initialise()
        if not started:
            print("FAIL. Unable to initialise all nodes.")
            runner.exit()
            exit(1)

        print("---INITIALISING TEST---")
        start = time()
        success = TestConsistencyAll(runner)
        if success:
            print(f"PASS. Time taken: {time()-start}")
        else:
            print(f"FAIL. Time taken: {time()-start}")

        print("Exiting.")
        runner.exit()
        
    except KeyboardInterrupt:
        print("Interrupted. Exiting.")
        if runner is not None:
            runner.exit()
    except Exception:
        print("Exception. Exiting.")
        if runner is not None:
            runner.exit()
        
