package seekret

import "fmt"

// print the entire server config
func (srv *Server) PrintServerConfig() {
	fmt.Println(srv.Config)
}
