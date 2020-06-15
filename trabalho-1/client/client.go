package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/V-H-R-Oliveira/Sistemas-Distribuidos/trabalho-1/data"
)

const serverEndpoint = "localhost:8080"

func sendPayload(ctx context.Context, conn net.Conn, payload []byte, recv chan<- data.Matricula) {
	conn.Write(payload)
	reader := json.NewDecoder(conn)

	for {
		var data data.Matricula
		err := reader.Decode(&data)

		if err != nil {
			return
		}

		select {
		case recv <- data:
		case <-ctx.Done():
			return
		}
	}
}

func client(ctx context.Context, matricula *data.Matricula) {
	payload, err := matricula.Serializar()

	if err != nil {
		log.Fatal(err)
	}

	recv := make(chan data.Matricula)
	defer close(recv)

	dial, err := net.Dial("tcp", serverEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer dial.Close()
	log.Println("[*] Conectado com o endpoint", dial.RemoteAddr().String())
	go sendPayload(ctx, dial, payload, recv)

	for {
		select {
		case data := <-recv:
			fmt.Printf("(%s) - %s -> %s\n", data.Ra, data.Aluno.Curso, data.Aluno.Curso)
			return
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	aluno := data.CriarAluno("João", "Biologia")
	matricula := data.CriarMatricula(aluno)

	client(ctx, matricula)
}
