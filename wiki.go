package main

import (
	"fmt"
	"io/ioutil"
)

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// salvar persiste os dados da página
func (p *Pagina) salvar() error {
	nomeDoArquivo := p.Titulo + ".txt"
	return ioutil.WriteFile(nomeDoArquivo, p.Corpo, 0600)
}

// Além de salvar páginas, também queremos carregar páginas:
func carregaPagina(titulo string) (*Pagina, error) {
	nomeDoArquivo := titulo + ".txt"
	corpo, err := ioutil.ReadFile(nomeDoArquivo)
	if err != nil {
		return nil, err
	}
	return &Pagina{Titulo: titulo, Corpo: corpo}, nil
}

func main() {
	escrever := &Pagina{Titulo: "PaginaTeste", Corpo: []byte("Esta é uma pagina de exemplo.")}
	escrever.salvar()
	ler, _ := carregaPagina("PaginaTeste")
	fmt.Println(string(ler.Corpo))
}
