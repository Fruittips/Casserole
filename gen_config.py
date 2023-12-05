import json
from pathlib import Path
from typing import Dict, Any
from dataclasses import dataclass

ENCODING = "utf-8"
START_NODE_ID = 3000

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


class ConfigEncoder(json.JSONEncoder):
    def default(self, o: Config) -> Dict[str, Any]:
        return o.__dict__


def gen_system(n: int = 5, consistencyLevel:str="QUORUM", gracePeriod=100, timeout=10, rf=3) -> Config:
    # Generate config
    nodesDict = {}
    for node_id in range(START_NODE_ID, START_NODE_ID + n):
        nodesDict[str(node_id)] = Node(node_id, False)
    config = Config(consistencyLevel, gracePeriod, timeout, rf, nodesDict)
    
    # Generate config.json file
    configFile = Path("./config.json")
    configFile.unlink(missing_ok=True)
    configFile.touch()
    configFile.write_text(config.to_json())

    # Generate new db file and hh files
    Path("./dbFiles").mkdir(exist_ok=True)
    Path("./hintedHandoffs").mkdir(exist_ok=True)
    for node_id in config.ring:
        dbFile = Path(f"./dbFiles/node-{node_id}.json")
        hhFile = Path(f"./hintedHandoffs/node-{node_id}.json")
        dbFile.unlink(missing_ok=True)
        hhFile.unlink(missing_ok=True)
        dbFile.touch()
        hhFile.touch()
        dbFile.write_text(json.dumps({
            "TableName": "Student Courses",
            "PartitionKey": "CourseId",
            "Partitions": {}
        }))
        hhFile.write_text("{}")

    __generate_dashboard(config)
    
    print(f"Generated system with {n} nodes, RF={rf}, consistency level {consistencyLevel}.")
    return config

def __generate_dashboard(config):
    html_template_builder = ""
    array_str = "["
    base_template="""
<div id="nodePORT_PLACEHOLDER" class="node bg-white shadow-md p-4 rounded w-1/3 min-w-xl w-full border-gray-300">
            <h3 class="text-2xl font-bold mb-6">Node PORT_PLACEHOLDER</h3>
            <div class="data-section mb-4">
                <strong class="block text-lg pb-1">Node Data:</strong>
                <p id="node-PORT_PLACEHOLDER-data">Status: Awaiting data...</p>
            </div>
            <div class="data-section mb-4">
                <strong class="block text-lg pb-1">DB Data:</strong>
                <p id="node-PORT_PLACEHOLDER-db-data">Status: Awaiting data...</p>
            </div>
            <div class="data-section mb-4">
                <strong class="block text-lg pb-1">HH Data:</strong>
                <p id="node-PORT_PLACEHOLDER-hh-data">Status: Awaiting data...</p>
            </div>

            <div class="w-full border-b border-gray-500 my-4"></div>

            <form onsubmit="postRequest(PORT_PLACEHOLDER); return false;"
                class="mb-8 border border-gray-300 rounded p-4 bg-gray-100">
                <label class="block mb-2">Course ID:
                    <input type="text" id="nodePORT_PLACEHOLDER-course-id" class="p-2 border rounded border-gray-300">
                </label>
                <label class="block mb-2">Student Name:
                    <input type="text" id="nodePORT_PLACEHOLDER-student-name" class="p-2 border rounded border-gray-300">
                </label>
                <label class="block mb-2">Student Number:
                    <input type="text" id="nodePORT_PLACEHOLDER-student-id" class="p-2 border rounded border-gray-300">
                </label>
                <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">
                    POST Data
                </button>
            </form>

            <form onsubmit="getRequest(PORT_PLACEHOLDER); return false;" class="border border-gray-300 rounded p-4 bg-gray-100">
                <label class="block mb-2">Course ID: <input type="text" id="nodePORT_PLACEHOLDER-course-id-get"
                        class="p-2 border rounded border-gray-300"></label>
                <label class="block mb-2">Student ID: <input type="text" id="nodePORT_PLACEHOLDER-student-id-get"
                        class="p-2 border rounded border-gray-300"></label>
                <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">
                    GET Data
                </button>
            </form>

            <button onclick="kill(PORT_PLACEHOLDER)"
                class="mt-8 bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">KILL</button>
            <button onclick="revive(PORT_PLACEHOLDER)"
                class="mt-8 bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700">REVIVE</button>
        </div>
"""
    for node_id in config.ring:
        port_html = base_template.replace("PORT_PLACEHOLDER", str(node_id))
        html_template_builder += port_html + "\n"
        array_str += f"{str(node_id)},"
    array_str += "]"
    with open("./frontend/base.html", 'r') as file:
        file_contents = file.read()
        file_contents = file_contents.replace("<!-- ===REPLACE NODE DASHBOARD=== -->", html_template_builder)
        file_contents = file_contents.replace("// ===REPLACE PORTS ARRAY===", f"const ports = {array_str};")
        with open("./frontend/casseroleChef.html", 'w') as file:
            file.write(file_contents)
    
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
    gen_system(
        args.nodeCount,
        args.consistencyLevel,
        100,
        args.timeout,
        args.rf,
    )
