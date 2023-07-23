package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ClienteDTO struct {
	Nome      string  `json:"nome"`
	Idade     uint8   `json:"idade"`
	TipoConta string  `json:"tipoConta"`
	Salario   float64 `json:"salario"`
}

type Cliente struct {
	Nome      string    `json:"nome"`
	Idade     uint8     `json:"idade"`
	TipoConta TipoConta `json:"tipoConta"`
	salario   float64   //`json:"salario"`
}

type TipoConta uint8

const (
	Cancelada TipoConta = iota
	BlackList
	PlanoBasico
	ContaSalario
	Premium
)

func (c *Cliente) parseJson(requestBody *io.ReadCloser) (e error) {

	clienteDTO := &ClienteDTO{}
	if e = json.NewDecoder(*requestBody).Decode(clienteDTO); e != nil {
		return
	}

	var tipoConta TipoConta
	switch clienteDTO.TipoConta {
	case Cancelada.String():
		tipoConta = Cancelada
	case BlackList.String():
		tipoConta = BlackList
	case ContaSalario.String():
		tipoConta = ContaSalario
	case Premium.String():
		tipoConta = Premium
	}
	*c = Cliente{Nome: clienteDTO.Nome, Idade: clienteDTO.Idade, TipoConta: tipoConta, salario: clienteDTO.Salario}
	return
}

func (t TipoConta) String() string {
	switch t {
	case Cancelada:
		return "Cancelada"
	case BlackList:
		return "BlackList"
	case ContaSalario:
		return "ContaSalario"
	case Premium:
		return "Premium"
	}
	return "unknown"
}

func main() {

	mux := http.NewServeMux()
	os.Setenv("PORT", ":8080")

	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/client", ClientHandler)
	mux.HandleFunc("/client/", ClientHandler)

	server := &http.Server{
		Addr:         "0.0.0.0" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Println("Server Start...")

	if e := server.ListenAndServe(); e != nil {
		log.Fatal(e.Error())
	}

}

func Healthz(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Server Status Run..."))
}

func ClientHandler(response http.ResponseWriter, request *http.Request) {

	fmt.Println(request.URL.RawFragment)
	fmt.Println(request.URL.Path)
	fmt.Println(request.URL.RawQuery)

	switch request.Method {
	case "GET":
		GetClient(response, request)
	case "POST":
		PostClient(response, request)
	case "DELETE":
		DeleteClient(response, request)
	case "PUT":
		PutClient(response, request)
	default:
		ErrorResponse(response, request)
	}
}

func PostClient(writer http.ResponseWriter, request *http.Request) {

	fmt.Println(request.URL)
	splitedUrl := strings.Split(request.URL.Path, "/")
	fmt.Println(splitedUrl)

	if len(request.URL.Path) > 0 && len(splitedUrl) != 2 || !strings.EqualFold("client", splitedUrl[1]) {
		request.Body.Close()
		writer.WriteHeader(http.StatusPreconditionFailed)
		writer.Write([]byte("No methodo post so é necessario enviar o body."))
		return
	}
	cliente := Cliente{}
	if e := cliente.parseJson(&request.Body); e != nil {
		request.Body.Close()
		writer.WriteHeader(http.StatusPreconditionFailed)
		writer.Write([]byte("Invalid BodyFormating."))
		return
	}
	fmt.Println(cliente)
	log.Println("Client posted...")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Client posted..."))
	request.Body.Close()
	return
}

func DeleteClient(writer http.ResponseWriter, request *http.Request) {
	log.Println("Client deleted...")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Client deleted..."))
}

func PutClient(writer http.ResponseWriter, request *http.Request) {
	log.Println("Client updated...")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Client updated..."))
}

func GetClient(writer http.ResponseWriter, request *http.Request) {
	log.Println("Client found...")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Client found..."))
}

func ErrorResponse(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusForbidden)
	response.Write([]byte("Method not allowed."))
}
