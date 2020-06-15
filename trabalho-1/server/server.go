package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/V-H-R-Oliveira/Sistemas-Distribuidos/trabalho-1/data"
	"github.com/twinj/uuid"
)

const serverAddr = ":8080"

type mapConn struct {
	conn      net.Conn
	matricula data.Matricula
}

func acceptConnections(ctx context.Context, listener net.Listener, pool chan<- net.Conn) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			return
		}

		defer conn.Close()
		log.Println("Received a connection from", conn.RemoteAddr().String())

		select {
		case pool <- conn:
		case <-ctx.Done():
			return
		}
	}
}

func recvData(ctx context.Context, conn net.Conn, recv chan<- mapConn, deadConns chan<- net.Conn) {
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
	case <-ctx.Done():
		return
	}
}

func server(ctx context.Context) {
	connPool := make(chan net.Conn)
	deadConnPool := make(chan net.Conn)
	recv := make(chan mapConn)

	listener, err := net.Listen("tcp", serverAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Println("Server listening at", listener.Addr().String())
	go acceptConnections(ctx, listener, connPool)

	for {
		select {
		case conn := <-connPool:
			go recvData(ctx, conn, recv, deadConnPool)
		case data := <-recv:
			data.matricula.Ra = uuid.NewV4().String()
			content, err := data.matricula.Serializar()

			if err == nil {
				data.conn.Write(content)
			}
		case conn := <-deadConnPool:
			log.Printf("Connection %s closed...\n", conn.RemoteAddr().String())
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server(ctx)
}
