# GNetCo - Go Netconf basic library
This is netconf library with minimal high level functions. This time it doesn't perform any error handling.
All you need to do is make new connection and send valid rpc XML request to the device.
TimeMe function is used to get rpc execution time.

### Quickstart
```golang
package main

import (
	"fmt"
	"gnetco/packme"
)

func main() {
	c, _ := packme.Connect("10.0.0.1:830", "UserName", "Password")
	defer c.Disconnect()

	// rpcGetRunningConfig := `<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><get-config><source><running/></source></get-config></rpc>`

	// elapsed := c.TimeMe(c.Exec, []byte(rpcGetRunningConfig))
	// log.Printf("GetRunningConfig took %s", elapsed)

	hn := `<rpc message-id="105" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0"><get-sessions/></rpc>`

	res, _ := c.Exec([]byte(hn))
	fmt.Println(string(res))
}
```

Copyright (c) 2022