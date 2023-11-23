import json
from random import choice
from os import chdir
from time import time, time_ns, sleep
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, read_req


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
    print("6. Sleeping for 10 seconds.")
    sleep(10)
    print(f"6. Read from node {svr_id}")
    resp, success = read_req(svr_id, course_id, "1006969")
    if not success:
        print(f"6. Error: {resp}")
        return False
    print(f"6. Read Response: {resp}")
    sleep(5)
    print(f"7. Read from node {svr_id}")
    resp, success = read_req(svr_id, course_id, "1006969")
    if not success:
        print(f"7. Error: {resp}")
        return False
    print(f"7. Read Response: {resp}")

    # Expect data to be the same
    # read_dat = json.loads(resp)
    # if read_dat["StudentId"] != data["StudentId"] or read_dat["StudentName"] != data["StudentName"]:
    #     print("2. Error: Data not matching")
    #     return False
    return True

if __name__ == "__main__":
    runner = None
    try:
        # Initialise project directory as root
        chdir("../")
        config = getConfigWithNNodes(10)
        runner = Runner(config)
        runner.initialise()

        print("---INITIALISING TEST---")
        start = time()
        success = TestMultiWriteOfDiffData(runner)
        if success:
            print(f"PASS. Time taken: {time()-start}")
        else:
            print(f"FAIL. Time taken: {time()-start}")

        print("Exiting.")
        #runner.exit()
        while True:
            pass
        
    except KeyboardInterrupt:
        print("Interrupted. Exiting.")
        if runner is not None:
            runner.exit()
