package core

import (
	"fmt"
	"mi/connectpool"
	"net"
)

// import "mi/"

func tcpServer() {
	listener, err := net.Listen("tcp", ":7898")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	// defer listener.Close()

	pool := connectpool.NewPool(4)
	go func(p *connectpool.Pool) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept is error!")
				return
			}
			go connectpool.Handle(conn, p)
		}
	}(pool)

	pool.Run()
}
