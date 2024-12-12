# WebSocket Server com Autenticação JWT

Este é um servidor WebSocket em Go que utiliza autenticação JWT (JSON Web Tokens) para permitir conexões seguras de clientes. Ele aceita conexões WebSocket apenas de origens confiáveis e valida o token JWT para garantir a autenticidade da solicitação.

## Funcionalidades

- **Autenticação com JWT**: Verifica o token JWT presente no cabeçalho `Authorization` da requisição.
- **Conexão WebSocket**: Estabelece uma conexão WebSocket com o cliente após autenticação bem-sucedida.
- **Envio de mensagens**: Envia e recebe mensagens entre clientes conectados.
- **Ping/Pong**: Utiliza o protocolo Ping/Pong do WebSocket para manter a conexão viva.
- **Requisições Seguras**: Limita conexões de WebSocket apenas a origens confiáveis.
- **Log de Atividades**: Registra eventos importantes e erros no arquivo `websocket_server.log`.

## Como usar

### Pré-requisitos

- [Go](https://golang.org/doc/install) 1.16+.
- Uma chave secreta para validar o JWT.

### Passos

1. Clone o repositório:

    ```bash
    git clone https://github.com/seu-usuario/seu-repositorio.git
    cd seu-repositorio
    ```

2. Instale as dependências:

    Execute o comando abaixo para instalar os pacotes necessários:

    ```bash
    go mod tidy
    ```

3. Atualize a chave secreta para o JWT:

    No arquivo `main.go`, altere a linha onde a chave secreta é definida:

    ```go
    return []byte("your-secret-key"), nil
    ```

    Substitua `"your-secret-key"` pela sua chave secreta.

4. Execute o servidor:

    ```bash
    go run main.go
    ```

    O servidor WebSocket estará disponível em `http://localhost:8080`.

5. Acesse a página WebSocket:

    Abra seu navegador e navegue até `http://localhost:8080`, onde você pode conectar-se ao servidor WebSocket.

6. Envie um token JWT válido:

    Ao conectar-se, envie um cabeçalho `Authorization` com o valor `Bearer <token>` para autenticar a conexão.

### Estrutura do Código

- **`main.go`**: Arquivo principal que contém a lógica do servidor WebSocket, autenticação JWT e manipulação de conexões.
- **`websocket_server.log`**: Arquivo de log gerado para registrar eventos do servidor.

### Arquitetura

1. O servidor inicia e escuta conexões HTTP na porta `8080`.
2. Quando um cliente tenta se conectar via WebSocket, o servidor verifica se o token JWT no cabeçalho de autorização é válido.
3. Caso o token seja válido, o servidor realiza o upgrade da conexão HTTP para WebSocket e permite que o cliente envie/receba mensagens.
4. O servidor envia periodicamente pings para manter as conexões vivas.
5. A cada nova mensagem recebida de um cliente, o servidor transmite essa mensagem para todos os outros clientes conectados.

## Log

Os logs do servidor são registrados no arquivo `websocket_server.log`, que pode ser consultado para diagnósticos e auditoria.

## Contribuindo

Se você quiser contribuir com melhorias ou correções para o projeto, basta fazer um fork deste repositório, realizar as modificações e enviar um pull request.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
