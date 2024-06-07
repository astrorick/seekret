package seekret

import "fmt"

// print the entire server config
func (s *Server) PrintServerConfig() {
	fmt.Println(s.ServerConfig)
}
