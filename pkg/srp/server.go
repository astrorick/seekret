package srp

type SRPServer struct {
	Username string
	Salt     []byte
	Verifier []byte
}

func (p *SRPParams) NewServer(username string, salt []byte, verifier []byte) (*SRPServer, error) {
	return nil, nil
}
