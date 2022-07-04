package packme

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

const enderString = "]]>]]>"

type GNetCo interface {
	connect(string, string, string) error
	Exec([]byte) ([]byte, error)
	Lock(func([]byte) ([]byte, error), []byte) ([]byte, error)
	unlock()
	Disconnect()
	TimeMe(func([]byte) ([]byte, error), []byte) time.Duration
}

type netconf struct {
	// Info   string
	reader  io.Reader
	writer  io.Writer
	session *ssh.Session
}

func Connect(socket, username, password string) (GNetCo, error) {
	n := new(netconf)
	n.connect(socket, username, password)
	// if err := n.connect(socket, username, password); err != nil {
	// }
	rpcHelloString := `<?xml version="1.0" encoding="UTF-8"?>
	<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
		<capabilities>
			<capability>urn:ietf:params:netconf:base:1.0</capability>
		</capabilities>
	</hello>`

	n.Exec([]byte(rpcHelloString))
	// if reply, err := n.Exec([]byte(rpcHelloString)); err != nil {
	// 	return nil, fmt.Errorf("%s", err)
	// 	// fmt.Println(string(reply))
	// } else {
	// 	fmt.Println(string(reply))
	// }
	return n, nil

}

func (n *netconf) connect(socket, username, password string) error {
	credentials := &ssh.ClientConfig{
		User:            username, //"netconf-user",
		Auth:            []ssh.AuthMethod{ssh.Password(password /*"netconf-password"*/)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", socket, credentials)
	if err != nil {
		return fmt.Errorf("ssh client dial %s", err)
	}

	n.session, err = sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("ssh client %s", err)
	}

	n.writer, err = n.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("writer error %s", err)
	}

	n.reader, err = n.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("reader error %s", err)
	}

	n.session.RequestSubsystem("netconf")

	return nil
}

func (n *netconf) Exec(rpcMesage []byte) ([]byte, error) {
	if err := n.send([]byte(rpcMesage)); err != nil {
		return nil, err
	}

	reply, err := n.receive()
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (n *netconf) TimeMe(f func([]byte) ([]byte, error), rpcMessage []byte) time.Duration {
	start := time.Now()
	f(rpcMessage)
	elapsed := time.Since(start)
	// log.Printf("GetRunningConfig took %s", elapsed)

	return elapsed
}

func (n *netconf) Lock(f func([]byte) ([]byte, error), rpcMessage []byte) ([]byte, error) {
	lockMessage := `<?xml version="1.0" encoding="UTF-8"?>
	<rpc message-id="101" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
		<lock>
		<target>
			<running/>
		</target>
		</lock>
	</rpc>`

	// start := time.Now()
	/*res, _ := */
	n.Exec([]byte(lockMessage))
	// elapsed := time.Since(start)
	// log.Printf("Lock took %s", elapsed)

	// fmt.Println(string(res))
	defer n.unlock()

	return f(rpcMessage)
}

func (n *netconf) unlock() {
	unlockMessage := `<?xml version="1.0" encoding="UTF-8"?>
	<rpc message-id="100" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
		<unlock>
		<target>
			<running/>
		</target>
		</unlock>
	</rpc>`
	n.Exec([]byte(unlockMessage))
}

func (n *netconf) Disconnect() {
	n.session.Close()
}

// func (n *netconf) send(rpcMessage string) error {
func (n *netconf) send(message []byte) error {
	message = append(message, []byte(enderString)...)

	_, err := n.writer.Write(message)
	if err != nil {
		return fmt.Errorf("writer error %s", err)
	}

	return nil
}

func (n *netconf) receive() ([]byte, error) {
	var receiveBytesBuffer bytes.Buffer
	buf := make([]byte, 8192)
	pos := 0

	for {
		n, err := n.reader.Read(buf[pos : pos+(len(buf)/2)])
		if err != nil {
			if err != io.EOF {
				log.Fatalf("reader error %s", err)
			}
			break
		}
		if n > 0 {
			end := bytes.Index(buf[0:pos+n], []byte(enderString))

			if end > -1 {
				receiveBytesBuffer.Write(buf[0:end])
				return receiveBytesBuffer.Bytes(), nil
				// break
			}

			if pos > 0 {
				receiveBytesBuffer.Write(buf[0:pos])
				copy(buf, buf[pos:pos+n])
			}

			pos = n
		}
	}
	return nil, fmt.Errorf("receive failed")
}
