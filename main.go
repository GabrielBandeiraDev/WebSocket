package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Apenas permitir conexões de origens confiáveis
		allowedOrigins := []string{"http://example.com", "https://example.com"}
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if strings.EqualFold(origin, allowedOrigin) {
				return true
			}
		}
		log.Println("Origin não permitida:", origin)
		return false
	},
}

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	name   string
	status string
	token  string
}

var clients = make(map[*Client]bool) // Mapa para armazenar clientes conectados
var mutex sync.Mutex                 // Mutex para garantir acesso concorrente seguro ao mapa

// Função que valida o token JWT
func validateToken(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}

	// Exemplo de validação de token (isso deve ser configurado de acordo com sua implementação de JWT)
	// Token = "Bearer <token>"
	token = strings.TrimPrefix(token, "Bearer ")
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Aqui você deve usar a chave secreta configurada para verificar o JWT
		return []byte("your-secret-key"), nil
	})

	return err == nil
}

// Função que lida com a conexão WebSocket
func handleConnection(w http.ResponseWriter, r *http.Request) {
	if !validateToken(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Atualiza a conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao upgrade da conexão:", err)
		return
	}
	client := &Client{
		conn:   conn,
		send:   make(chan []byte),
		name:   r.URL.Query().Get("name"),
		status: "online",
	}

	// Utilizar mutex para garantir thread-safe na manipulação de clientes
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	log.Printf("Novo cliente conectado: %s", client.name)

	// Definir tempo limite de leitura
	conn.SetReadDeadline(time.Now().Add(time.Second * 60))

	// Configurar Ping/Pong para manter a conexão viva
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		return nil
	})

	// Fechar a conexão de forma segura
	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		conn.Close()
		log.Printf("Cliente desconectado: %s", client.name)
	}()

	go readMessages(client)

	// Enviar mensagens para todos os clientes
	for msg := range client.send {
		mutex.Lock()
		for c := range clients {
			if c != client {
				err := c.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Erro ao enviar mensagem para", c.name, err)
					c.conn.Close()
					delete(clients, c)
				}
			}
		}
		mutex.Unlock()
	}
}

// Função que lê as mensagens enviadas pelo cliente
func readMessages(client *Client) {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Erro na leitura de mensagem:", err)
			return
		}
		log.Printf("Mensagem recebida de %s: %s", client.name, message)
		client.send <- message
	}
}

// Função de log para registrar mensagens de erro e eventos importantes
func setupLogging() {
	// Criar arquivo de log, se não existir
	logFile, err := os.OpenFile("websocket_server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Não foi possível abrir o arquivo de log:", err)
	}
	log.SetOutput(logFile)
	log.Println("Log iniciado")
}

// Função que trata reconexões automáticas
func handleReconnection(client *Client) {
	// Tentar reconectar em intervalos definidos
	retryInterval := 5 * time.Second
	for {
		time.Sleep(retryInterval)
		// Tente reconectar...
		// Este é um exemplo básico, você pode querer melhorar com backoff exponencial
	}
}

// Função para enviar uma mensagem de "ping" de vez em quando
func sendPingMessages() {
	for {
		mutex.Lock()
		for c := range clients {
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("Erro ao enviar ping para", c.name, err)
			}
		}
		mutex.Unlock()
		time.Sleep(30 * time.Second)
	}
}

// Função principal
func main() {
	// Inicializa o log
	setupLogging()

	// Configura a rota WebSocket
	http.HandleFunc("/ws", handleConnection)

	// Serve o arquivo HTML para os usuários
	http.Handle("/", http.FileServer(http.Dir("./")))

	// Iniciar a rotina de envio de ping
	go sendPingMessages()

	addr := "localhost:8080"
	fmt.Printf("Servidor WebSocket rodando em %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
