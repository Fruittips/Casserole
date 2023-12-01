import json
from random import choice
from os import chdir
from time import time, sleep
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, read_req
from sys import exit

def TestMultiWriteOfDiffData(runner: Runner) -> bool:
    """
    Tests with multiple writes of different data to the same node. The latest write SHOULD win (eventually).

    We do this by changing the name, since the things that are overwritten are:
    - Name
    - Timestamp Written
    - Timestamp Deleted
    """

    # Select random node to write to
    svr_id = choice(list(runner.config.ring.keys()))

    # Test Data
    course_id = "50.069"
    data1 = get_write_data("1006969", "TestUser1", 0, -1)
    data2 = get_write_data("1006969", "TestUser2", 0, -1)
    data3 = get_write_data("1006969", "TestUser3", 0, -1)
    data4 = get_write_data("1006969", "TestUser4", 0, -1)
    data5 = get_write_data("1006969", "TestUser5", 0, -1)

    # Write to node multiple times
    print(f"1. Write to node {svr_id} with data {data1}.")
    resp, success = write_req(svr_id, course_id, data1)
    if not success:
        print(f"1. Error: {resp}")
        return False
    print(f"1. Write Response: {resp}")

    sleep(1)
    
    print(f"2. Write to node {svr_id} with data {data2}.")
    resp, success = write_req(svr_id, course_id, data2)
    if not success:
        print(f"2. Error: {resp}")
        return False
    print(f"2. Write Response: {resp}")

    sleep(1)
    
    print(f"3. Write to node {svr_id} with data {data3}.")
    resp, success = write_req(svr_id, course_id, data3)
    if not success:
        print(f"3. Error: {resp}")
        return False
    print(f"3. Write Response: {resp}")

    sleep(1)
    
    print(f"4. Write to node {svr_id} with data {data4}.")
    resp, success = write_req(svr_id, course_id, data4)
    if not success:
        print(f"4. Error: {resp}")
        return False
    print(f"4. Write Response: {resp}")

    sleep(1)
    
    print(f"5. Write to node {svr_id} with data {data5}.")
    resp, success = write_req(svr_id, course_id, data5)
    if not success:
        print(f"5. Error: {resp}")
        return False
    print(f"5. Write Response: {resp}")

    # Read from node
    print("6. Sleeping for 5 seconds.")
    sleep(5)
    print(f"6. Read from node {svr_id}")
    resp, success = read_req(svr_id, course_id, "1006969")
    print(f"6. Read Response: {resp}")

    ## Determine if data is the expected data
    read_dat = json.loads(resp)
    success = read_dat["StudentName"] == "TestUser5" # We expect it to be the LATEST WRITTEN user.
    if not success:
        print("7. Record read was NOT the latest. Retrying after 5 seconds to account for read-repair...")
        sleep(5)
        print(f"7. Read from node {svr_id}")
        resp, success = read_req(svr_id, course_id, "1006969")
        print(f"7. Read Response: {resp}")
        read_dat = json.loads(resp)
        success = read_dat["StudentName"] == "TestUser5"

        if not success:
            print("Error: Record read was NOT the latest after two tries.")
            return False
    return True

if __name__ == "__main__":
    runner = None
    try:
        # Initialise project directory as root
        chdir("../")
        config = getConfigWithNNodes(10)
        runner = Runner(config)
        started = runner.initialise()
        if not started:
            print("FAIL. Unable to initialise all nodes.")
            runner.exit()
            exit(1)

        print("---INITIALISING TEST---")
        start = time()
        success = TestMultiWriteOfDiffData(runner)
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
