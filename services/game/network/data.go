package network

type ClientData struct {
	UID   uint64
	Name  string
	Token string
}

func (c *ClientData) verifyToken() error {
	// todo
	return nil
}
