package handlers

import (
	"encoding/json"
	"fmt"
	"os"
)

// still empty
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func internalWrite(curNode int, a AtomicDbMessage, toNode int) error {
	// check if config.json is_dead is false
	is_dead, err := checkIsDead()
	// write into json file service
	if is_dead {
		return err
	}

	d, err2 := internalRead(toNode)
	if err2 != nil {
		return err2
	}

	// Add the data into the struct
	strDeadNodeId := fmt.Sprintf("%d", toNode)
	d.Row[strDeadNodeId] = a

	// convert struct into bytes
	byteFile, err3 := json.MarshalIndent(d, "", "\t")
	if err3 != nil {
		return err3
	}

	// write byte into file
	filename := fmt.Sprintf("dbFiles/node-%d.json", toNode)
	err4 := os.WriteFile(filename, byteFile, 0644)
	if err3 != nil {
		return err4
	}

	return nil

}
