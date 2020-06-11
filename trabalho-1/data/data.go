package data

import "encoding/json"

// Aluno -> representa um determinado aluno
type Aluno struct {
	Nome  string `json:"nome"`
	Curso string `json:"curso"`
}

// Matricula -> Representa uma determinada matrícula
type Matricula struct {
	Ra    string `json:"ra"`
	Aluno *Aluno
}

// Operations -> Representa as operações possíveis
type Operations interface {
	Serializar() ([]byte, error)
}

// CriarAluno -> Cria um novo aluno
func CriarAluno(nome, curso string) *Aluno {
	return &Aluno{
		Nome:  nome,
		Curso: curso,
	}
}

// CriarMatricula -> Cria uma nova matrícula
func CriarMatricula(aluno *Aluno) *Matricula {
	return &Matricula{
		Ra:    "",
		Aluno: aluno,
	}
}

// Serializar -> Serializa uma matrícula no formato JSON
func (matricula *Matricula) Serializar() ([]byte, error) {
	return json.Marshal(matricula)
}
