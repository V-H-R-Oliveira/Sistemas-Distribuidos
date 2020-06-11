package main

import (
	"encoding/json"
	"fmt"
	"github.com/V-H-R-Oliveira/Sistemas-Distribuidos/trabalho-1/data"
	"log"
	"net"
)

const serverEndpoint = "localhost:8080"

func sendPayload(conn net.Conn, payload []byte, recv chan<- data.Matricula, cancel <-chan struct{}) {
	conn.Write(payload)
	reader := json.NewDecoder(conn)

	for {
		var data data.Matricula
		err := reader.Decode(&data)

		if err != nil {
			break
		}

		recv <- data
	}

	select {
	case <-cancel:
		return
	}
}

func client(matricula *data.Matricula, cancel <-chan struct{}) {
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
	go sendPayload(dial, payload, recv, cancel)

	for {
		select {
		case data := <-recv:
			fmt.Printf("(%s) - %s -> %s\n", data.Ra, data.Aluno.Curso, data.Aluno.Curso)
			return
		case <-cancel:
			return
		}
	}
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	aluno := data.CriarAluno("João", "Biologia")
	matricula := data.CriarMatricula(aluno)

	client(matricula, cancel)
}
