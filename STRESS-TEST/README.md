# STRESSER

É um CLI (Command Line Interface) desenhado para fazer GET requests a uma determinada URL, a fim de testar a carga suportada pela mesma.
Ele possui 3 flags para auxiliar:
- `url`: Aqui informe a URL que deseja se testar
- `requests`: Aqui informe a quantidade total de requests a serem enviadas.
- `concurrency`: Aqui informe a qunatidade de requests que deverá ser executada em concorrência.

## Como usar

Pode ser realizado a instalação desse CLI e utilizado da seguinte forma:
```sh
go install .
stresser --url https://...
```

Ou de preferir, utilizar um Container Docker (Publicado no DockerHub):
```sh
docker run devlucca/stresser --url https://...
```

### Flags

Aqui vai um exemplo do uso das Flags.

#### Exemplo 1:

Foi enviado o seguinte comando:
    
```sh
stresser --url https://google.com --requests 100 --concurrency = 5
```

O CLI irá executar 100 requests, de 5 em 5, sendo que a cada 5, estas estarão sendo executadas em paralelo. Assim que essas 5 requests terminarem, o CLI segue para as próximas 5.

#### Exemplo 2:

Foi enviado o seguinte comando:
    
```sh
stresser --url https://google.com --requests 8 --concurrency = 5
```

O CLI irá executar 8 requests, de 5 em 5, sendo que a cada 5, estas estarão sendo executadas em paralelo. Como depois da primeira execução, so restarão 3 requests, o executor vai executar as últimas 3 concorrentemente.


## Resultados

O resultado será um report como abaixo:
```
Stress Test - Request report.
=================================
Took 20.194663178s to complete
Total requests:       100 requests
Completed requests:   100 requests

---------------------------------
Time Distribution:
    - P99th: 2.741275976s
    - P75th: 2.196581368s
    - P50th: 1.12500338s
---------------------------------
HTTP Status Code Distribution:
    - 2xx:             100 requests
    - 4xx:             0 requests
    - 5xx:             0 requests
---------------------------------
Failed requests:
    0 requests failed
```

Onde será possível visualizar algumas informações que foram coletadas.
