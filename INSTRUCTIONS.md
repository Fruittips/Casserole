# Managers
### Node Manager
```go
	node := utils.NewNodeManager(3001)

	// append to hinted handoffs
	node.HintedHandoffManager.Append(1, utils.AtomicDbMessage{Data: []string{"hello", "world"}, Timestamp: 123})

	// append to db
	node.DatabaseManager.AppendRow(utils.Row{Data: []string{"hello", "asdads"}, Timestamp: 123})

	// read from db
	fmt.Println(node.DatabaseManager.Data)

  // get row by id
  data, err := node.DatabaseManager.GetRowById(1)
  fmt.Println(data)

	// read from hinted handoffs
	fmt.Println(node.HintedHandoffManager.Data)

	// get config data
	fmt.Println(node.ConfigManager.Data)

	// get my node data
	fmt.Println(node.Me())
```




