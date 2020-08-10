package main

import (
	"fmt"
	"io/ioutil"

	// Para usar o net/http package, ele deve ser importado:
	"net/http"
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

// Vamos criar um Handler, viewHandler que permitirá aos usuários visualizar uma página wiki.
// Ele irá lidar com URLs prefixados com "/ view /".

// Primeiro, essa função extrai o título da página ler.URL.Path, o componente do caminho da URL requisitada.
// Path é dividido [len("/view/"):] para eliminar "/view/" componente principal do caminho da solicitação.
// Isso ocorre porque o caminho invariavelmente começará com "/view/", que não faz parte do título da página.
// carregaPagina então carrega os dados da página, formata a página com uma string de HTML simples
// e a grava no "escrever" http.ResponseWriter.

// Observe o uso de _ para ignorar o error value do retorno de carregaPagina. Isso é feito aqui para
// simplificar e geralmente é considerado uma prática ruim. Cuidaremos disso mais tarde.

func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/view/"):]
	pagina, _ := carregaPagina(titulo)
	fmt.Fprintf(escrever, "<h1>%s</h1><div>%s</div>", pagina.Titulo, pagina.Corpo)
}

func main() {

}
