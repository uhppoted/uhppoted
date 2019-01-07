package uhppote

import (
	"errors"
	"fmt"
	"net"
)

func close(connection net.Conn) {
	fmt.Println(" ... closing connection")

	connection.Close()
}

func makeErr(msg string, err error) error {
	return errors.New(fmt.Sprintf(msg+" [%v]", err))
}
