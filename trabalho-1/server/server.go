package main

import (
	"encoding/json"
	"log"
	"net"
	"tarefa/m/data"

	"github.com/twinj/uuid"
)

const serverAddr = ":8080"

type mapConn struct {
	conn      net.Conn
	matricula data.Matricula
}

func acceptConnections(listener net.Listener, pool chan<- net.Conn, cancel <-chan struct{}) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			return
		}

		defer conn.Close()
		log.Println("Received a connection from", conn.RemoteAddr().String())

		select {
		case pool <- conn:
		case <-cancel:
			return
		}
	}
}

func recvData(conn net.Conn, recv chan<- mapConn, deadConns chan<- net.Conn, cancel <-chan struct{}) {
	reader := json.NewDecoder(conn)

	for {
		var data data.Matricula
		err := reader.Decode(&data)

		if err != nil {
			break
		}

		recv <- mapConn{
			conn:      conn,
			matricula: data,
		}
	}

	select {
	case deadConns <- conn:
	case <-cancel:
		return
	}
}

func server(cancel <-chan struct{}) {
	connPool := make(chan net.Conn)
	defer close(connPool)

	deadConnPool := make(chan net.Conn)
	defer close(deadConnPool)

	recv := make(chan mapConn)
	defer close(recv)

	listener, err := net.Listen("tcp", serverAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Println("Server listening at", listener.Addr().String())
	go acceptConnections(listener, connPool, cancel)

	for {
		select {
		case conn := <-connPool:
			go recvData(conn, recv, deadConnPool, cancel)
		case data := <-recv:
			data.matricula.Ra = uuid.NewV4().String()
			content, err := data.matricula.Serializar()

			if err == nil {
				data.conn.Write(content)
			}
		case conn := <-deadConnPool:
			log.Printf("Connection %s closed...\n", conn.RemoteAddr().String())
		case <-cancel:
			return
		}
	}
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	server(cancel)
}
