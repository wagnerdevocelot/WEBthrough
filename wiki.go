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

// Neste ponto, temos uma estrutura de dados simples e a capacidade de salvar e carregar um arquivo.
// Vamos escrever uma main func para testar o que escrevemos:

func main() {
	escrever := &Pagina{Titulo: "PaginaTeste", Corpo: []byte("Esta é uma pagina de exemplo.")}
	escrever.salvar()
	ler, _ := carregaPagina("Pagina teste")
	fmt.Println(string(ler.Corpo))
}

// Após compilar e executar este código, um arquivo chamado PaginaTeste.txt será criado,
// contendo o conteúdo de escrever. O arquivo será então lido na estrutura ler e seu corpo impresso na tela.

// Você pode compilar e executar o programa assim:

// $ go build wiki.go
// $ ./wiki
// Esta é uma pagina de exemplo.
