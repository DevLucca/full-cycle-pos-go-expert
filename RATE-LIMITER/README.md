# RATE LIMITER

Foi criado uma biblioteca, na qual subi em um outro [repositório](https://github.com/lccmrx/gorlim) e chamei carinhosamente de Gorlim (Go Rate Limiter).

é possível configurar o Gorlim de diversas formas, programaticamente (runtime) ou environments. No [Compose](./compose.yaml#api), deixei configurado uns valores simbólicos.

## Para subir

Criei uma API simples com o Go HTTP. E associei o Gorlim a ela.
```sh
docker compose up --build -d
```

Subirá todas as dependencias necessárias para realizar os testes.

Após a subida dos containers, pode se executar as chamadas HTTP, como exemplo:
```curl
curl localhost:8080/
```
ou
```curl
curl localhost:8080/ -H 'x-api-key: abc'
```

## Disclaimer

Como estou utilizando alguns Headers para busca de IP e alguns fallbacks, as vezes o IP retornado contém a porta na qual saiu a request. Isto em um servidor remoto não seria um problema, no entando, ao testar localmente, acaba incomodando os testes.

Recomendo que seja incluído o Header `X-Real-Ip` ou `X-Forwarded-For` com um valor fixo para melhores resultados!
