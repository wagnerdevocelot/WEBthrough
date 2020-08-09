package main

// Vamos começar definindo as estruturas de dados. Um wiki consiste em uma série de páginas interconectadas,
// cada uma com um título e um corpo (o conteúdo da página). Aqui, definimos Pagina como uma estrutura com dois
// campos representando o título e o corpo.

// Pagina é um struct com a estrutura que as páginas da aplicação terão
type Pagina struct {
	Titulo string
	Corpo  []byte
}

// O tipo []byte significa "um byte slice".
// Consulte slices https://golang.org/doc/articles/slices_usage_and_internals.html para obter mais informações
// O elemento Corpo é um []byte e não string porque esse é o tipo esperado pelos pacotes de io que usaremos.

func main() {

}
