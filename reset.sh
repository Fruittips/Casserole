#!/bin/bash

# This util script was generated using ChatGPT

# Array of port numbers
ports=(3000 3001 3002 3003 3004)

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
done

# Remove the trailing comma and close the JSON
config_content="${config_content%,}
    }   
}"

# Remove the trailing "& " from run_commands
run_commands="${run_commands%& }"

echo "$config_content" > ./config.json

echo "$run_commands" > ./run.sh

echo "Files created successfully!"
