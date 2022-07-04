# GNetCo - Go Netconf basic library
Simple library  which works with

```go
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

	// elapsed = c.TimeMe(c.Exec, []byte(hn))
	// log.Printf("GetHN took %s", elapsed)

	res, _ := c.Exec([]byte(hn))
	fmt.Println(string(res))
}

```
