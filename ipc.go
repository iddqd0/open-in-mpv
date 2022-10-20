package open_in_mpv

import (
	"net"
	"time"
	"runtime"
	"fmt"
	
	npipe "gopkg.in/natefinch/npipe.v2"
)

var ipcConnection Ipc

// Defines an IPC connection with a UNIX socket
type Ipc struct {
	conn          net.Conn
	SocketAddress string
}

// Send a byte-encoded command to the specified UNIX socket
func (i *Ipc) Send(cmd []byte) error {
	var err error
	
	switch runtime.GOOS {
		case "windows":
			timeout := time.Millisecond * 150
			i.conn, err = npipe.DialTimeout(`\\.\pipe\` + i.SocketAddress, timeout)
		default:
			i.conn, err = net.Dial("unix", i.SocketAddress)
	}
	
	if err != nil {
		return err
	}
	defer i.conn.Close()

	// The command has to be newline terminated
	if cmd[len(cmd)-1] != '\n' {
		cmd = append(cmd, '\n')
	}

	if _, err := fmt.Fprintln(i.conn, string(cmd)); err != nil {
		return err
	}

	return nil
}

// Generic public send string command for the default connection
func SendString(cmd string) error {
	return ipcConnection.Send([]byte(cmd))
}

// Generic public send byte-encoded string command for the default connection
func SendBytes(cmd []byte) error {
	return ipcConnection.Send(cmd)
}

func IpcConnect(path string) error {
	ipcConnection = Ipc{
		SocketAddress: path,
	}
	return nil
}