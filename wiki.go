package main

// importação do pacote io/ioutil
import "io/ioutil"

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// A struct Pagina descreve como os dados da página serão armazenados na memória.
// Mas e o armazenamento persistente? Podemos resolver isso criando um método salvar:

func (p *Pagina) salvar() error {
	nomeDoArquivo := p.Titulo + ".txt"
	return ioutil.WriteFile(nomeDoArquivo, p.Corpo, 0600)
}

// O que o método salvar faz? Ele tem um receiver que aponta para Pagina e tem um argumento "p"
// salvar não possui parametros além do presente no receiver, e retorna um error.

// nomeDoArquivo é uma varável que armazena o valor, Titulo + a string ".txt" que é a extenção

// O retorno de salvar é um error pois esse é o mesmo tipo que WriteFile retorna.
// O que WriteFile faz? WriteFile escreve um slice de bytes em um arquivo se salvar funcionar sem erros
// ele retorna nil que é o zero value para pointers.

// O literal inteiro 0600, passado como o terceiro parâmetro para WriteFile, indica que o arquivo
// deve ser criado com permissões de leitura e gravação apenas para o usuário atual.
// (Veja a página de manual do Unix open(2) para mais detalhes.)

func main() {

}
