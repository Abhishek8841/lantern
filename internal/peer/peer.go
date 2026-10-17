
package peer

import "net"

type Peer struct {
	ID   string
	Conn net.Conn
}
