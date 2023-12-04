from os import chdir
from time import time, sleep
from test_driver import Runner, get_write_data, getConfigWithNNodes, write_req, kill_req, revive_req
from sys import exit

def TestHH(runner: Runner) -> bool:
    """
    Triggers a sequence of events that causes a hintedhandoff to occur.
    """

    # Select random nodes
    write_node_id = "3000"
    kill_node_id = "3001"

    # Test Data
    course_id = "50.069"
    student_id = "1006969"
    student_name = "TestUser1"

    print("\n\n---SIMULATING HINTED HANDOFF.---")
    print("We use a system of 3 nodes with RF=3 and consistency QUORUM.")
    print("We kill a node, write to another, then revive the node.")
    print("This should trigger a hinted handoff to the once-dead node.")

    input("ENTER to kill node...")
    print()

    # Kill Node
    print(f"0. Kill N{kill_node_id}")
    success = kill_req(kill_node_id)
    if not success:
        print(f"0. Error: Could not kill N{kill_node_id}")
        return False
    print(f"0. Killed N{kill_node_id}.")

    input("ENTER to write to node...")
    print()

    # Write to Node
    data = get_write_data(student_id, student_name, 0, -1)
    print(f"1. Write to node {write_node_id} with data {data}.")
    resp, success = write_req(write_node_id, course_id, data)
    if not success:
        print(f"1. Error: {resp}")
        return False
    print(f"1. Write Response: {resp}")

    input("ENTER to revive node...")
    print()

    # Revive node, should trigger HH
    success = revive_req(kill_node_id)
    if not success:
        print(f"2. Error: Could not revive N{kill_node_id}")
        return False
    print(f"2. Revived N{kill_node_id}. This should trigger a HH.")

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
        success = TestHH(runner)
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
