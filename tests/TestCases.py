import requests
import json
from random import choice
from os import chdir
from typing import Dict, Any, Callable
from time import time, time_ns
from test_driver import Config, Node, Runner, get_read_url, get_write_url, get_write_data, getConfigWithNNodes

def read_req(svr_id: str, course_id: str, student_id: str) -> (str, bool):
    try:
        resp = requests.get(get_read_url(svr_id, course_id, student_id))
        if resp.status_code != 200:
            return f"HTTP Error {resp.status_code}, resp.text", False
        return resp.text, True
    except requests.exceptions.Timeout:
        return f"Timeout", False
    except requests.exceptions.TooManyRedirects:
        return f"Too Many Redirects", False
    except requests.exceptions.RequestException as e:
        return f"Catastrophic Error {e}", False

def write_req(svr_id: str, course_id: str, data: Dict[str, Any]) -> (str, bool):
    try:
        resp = requests.post(get_write_url(svr_id, course_id), data=data)
        if resp.status_code != 200:
            return f"HTTP Error {resp.status_code}, resp.text", False
        return resp.text, True
    except requests.exceptions.Timeout:
        return f"Timeout", False
    except requests.exceptions.TooManyRedirects:
        return f"Too Many Redirects", False
    except requests.exceptions.RequestException as e:
        return f"Catastrophic Error {e}", False

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
        print(f"2. Error: Data not matching")
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
        
