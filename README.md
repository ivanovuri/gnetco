# GNetCo - Go NETCONF basic library

This is NETCONF library with minimal high level functions. This time it doesn't perform any error handling, but it will do.
All you need to do is make new connection and send valid rpc XML request to the device.

Lock function is used to lock configuration during performing Exec function, where unlock performed after exiting Lock function.
TimeMe function is used to get rpc execution time. Just for understanding how long process take place.

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