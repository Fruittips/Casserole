#!/bin/bash

# This util script was generated using ChatGPT

# Array of port numbers
ports=(3000 3001 3002 3003)

# Ensure the directories exist, if not create them
mkdir -p ./dbFiles
mkdir -p ./hintedHandoffs

config_content=$(cat <<EOF
{
    "consistencyLevel": "QUORUM",
    "gracePeriod": 100,
    "timeout": 10,
    "rf": 3,
    "ring": {
EOF
)
run_commands=""

node_array=""
node_template=$(cat <<- EOM
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
EOM
)


# Loop through each port number
for port in "${ports[@]}"; do

  # Create the JSON content for dbFiles
  dbfile_content=$(cat <<EOF
{
  "TableName": "Student Courses",
  "PartitionKey": "CourseId",
  "Partitions": {
  }
}
EOF
  )

  # Create the JSON content for hintedHandoffs
  hinted_content=$(cat <<EOF
{
  "Row": {
  }
}
EOF
  )

  # Write the content to a file in /dbFiles named node-portNumber.json
  echo "$dbfile_content" > "./dbFiles/node-$port.json"

  # Write the content to a file in /hintedHandoffs named node-port.json
  echo "$hinted_content" > "./hintedHandoffs/node-$port.json"

  config_content+=$(cat <<EOF
        "$port": {
            "port": $port,
            "isDead": false
        },
EOF
    )
  
  run_commands+="go run . -port=$port & "

  node_with_port=${node_template//PORT_PLACEHOLDER/$port}
  node_array+="$node_with_port"
done

# Remove the trailing comma and close the JSON
config_content="${config_content%,}
    }   
}"

# Remove the trailing "& " from run_commands
run_commands="${run_commands%& }"

echo "$config_content" > ./config.json

echo "$run_commands" > ./run.sh

cp ./frontend/base.txt ./frontend/casseroleChef.html
# Replace in the casseroleChef.html file
perl -pi -e "s|===REPLACE NODE DASHBOARD===|$node_array|gs" frontend/casseroleChef.html

# Replace the ports array
IFS=,  # Set the internal field separator to comma
ports_array="const ports = [${ports[*]}];"
sed -i "" "s|===REPLACE PORTS ARRAY===|$ports_array|g" frontend/casseroleChef.html
IFS=   # Reset the internal field separator

echo "Files created successfully!"
