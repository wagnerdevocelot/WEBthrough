package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

// viewHandler escreve o titulo e corpo da pagina em html formatado
func viewHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/view/"):]
	pagina, _ := carregaPagina(titulo)
	fmt.Fprintf(escrever, "<h1>%s</h1><div>%s</div>", pagina.Titulo, pagina.Corpo)
}

// A função editHandler carrega a página (se não existir, cria uma Pagina struct vazia )
// e exibe um formulário HTML
// Esta função funcionará bem, mas todo esse HTML hardcoded é feio. Existe uma maneira melhor.

func editHandler(escrever http.ResponseWriter, ler *http.Request) {
	titulo := ler.URL.Path[len("/edit/"):]
	pagina, err := carregaPagina(titulo)
	if err != nil {
		pagina = &Pagina{Titulo: titulo}
	}
	fmt.Fprintf(escrever, "<h1>Editando %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		pagina.Titulo, pagina.Titulo, pagina.Corpo)
}

// Um wiki não é um wiki sem a capacidade de editar páginas. Vamos criar dois novos manipuladores: um
// nomeado editHandler para exibir um formulário de 'página de edição' e outro nomeado saveHandler para salvar
// os dados inseridos por meio do formulário.

// OBS. saveHandler está comentado pois essa função ainda não foi feita, o Go irá acusar erro para
// variáveis não utilizadas ou chamadas de funções que não foram implementadas.

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
