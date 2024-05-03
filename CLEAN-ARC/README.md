# Instrucoes para executar

1. docker compose up -d
2. Aguarde uns 10 segundinhos para dar tempo do RMQ e Banco subirem
2. make run (utiliza o `go run`) | make run-bin (gera o binario e executa)

# Portas disponiveis
|Protocolo | Porta |
|----------|-------|
|HTTP      | 8000  |
|gRPC      | 50051 |
|GraphQL   | 8080  |