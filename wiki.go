package main

import "io/ioutil"

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

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

// Como carregaPagina lê nosso arquivo? carregaPagina recebe uma string como parâmetro e nessa vai o titulo
// do arquivo, carregaPagina retorna Pagina.

// nomeDoArquivo pega o nome passado em parametro como titulo e concatena com ".txt"

// ReadFile lê o arquivo. ReadFile tem mais de um retorno, um []byte e error, nesse caso o []byte é o corpo
// da pagina. O if statement logo após é um tratamento de erros comum em Go, nesse caso se tentarmos ler
// um arquivo que não existe a função retornará uma mensagem indicando o problema e pra que isso aconteça
// err precisa ser diferente de nil

func main() {

}
