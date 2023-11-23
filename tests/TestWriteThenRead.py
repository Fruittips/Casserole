import json
from random import choice
from os import chdir
from time import time, time_ns
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, read_req


def TestWriteThenRead(runner: Runner) -> bool:
    """ Tests a standard write and read """

    # Select random node to write to
    svr_id = choice(list(runner.config.ring.keys()))

    # Test Data
    course_id = "50.069"
    data = get_write_data("1006969", "TestUser", time_ns(), -1)

    # Write to node
    print(f"1. Write to node {svr_id} with data {data}.")
    resp, success = write_req(svr_id, course_id, data)
    if not success:
        print(f"1. Error: {resp}")
        return False
    print(f"1. Write Response: {resp}")

    # Read from node
    print(f"2. Read from node {svr_id}")
    resp, success = read_req(svr_id, course_id, data["StudentId"])
    if not success:
        print(f"2. Error: {resp}")
        return False
    print(f"2. Read Response: {resp}")

    # Expect data to be the same
    read_dat = json.loads(resp)
    if read_dat["StudentId"] != data["StudentId"] or read_dat["StudentName"] != data["StudentName"]:
        print("2. Error: Data not matching")
        return False
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
        success = TestWriteThenRead(runner)
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
        
