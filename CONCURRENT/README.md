# Fastest CEP API on all Western | API de CEP mais rápida de todo oeste

## How it works | Como funciona?
O fluxo iniciante da API se preocupa em sincronizar os starts de todas as chamadas de _Go Routines_, tentando ao máximo ser fidedigno ao tempo de início. Caso isso não fosse feito, haveria brecha para uma API ter alguns nanosegundos de `headstart` (que seria quase que imperceptível, porém, injusto.)

Abaixo um teste realizado com o `starter` _channel_ e _waitGroups_:
```plain text
calling `viacep`        -----    22:21:18.327552000
calling `brasilapi`     -----    22:21:18.327564000
```

Produzindo uma diferença de tempo entre a primeira e segunda de `12000` nanosegundos

e outro teste realizado sem o `starter` _channel_ e _waitGroups_:
```plain text
calling `brasilapi`     -----    22:22:41.692886000
calling `viacep`        -----    22:22:41.692930000
```

Produzindo uma diferença de tempo entre a primeira e segunda de `44000` nanosegundos

Então, pode ser dizer que a aproximação com o `starter` _channel_ acaba sendo mais fidedigna.


Com isso, realizo as chamadas de ambas APIs, e espero o retorno *VÁLIDO* da mais rápida.
Considero como válida a API que retornar com um _status code_ igual a 200 e que não tenha elevado nenhum tipo de erro no _runtime_.

Assim que recebo a resposta, utilizo do `response` _channel_ para capturar a resposta!

Ao mesmo tempo, a _main thread_ está sendo executada e contabilizando o tempo de 1 segundo, e caso nenhuma das APIs retorne em tempo hábil, a aplicação finaliza.

## How to build | Como buildar

Use o seguinte comando para gerar o executável/binário:
```sh
go build -o fast-cep
```

## How to use | Como utilizar

Estamos utilizando uma aproximação de inputs via argumentos. Então simplesmente informe o seu CEP juntamente com a inicialização da aplicação, e tudo deve correr bem!

Formato `go run`
```sh
go run main.go 01234001
```
ou

Formato `binário` (depende do build do seu sistema operacional)
```sh
./fast-cep 01234001
```



