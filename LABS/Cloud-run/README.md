# Fluxo
A aplicação busca o CEP na ViaCEP, e a partir do retorno, utilizamos a cidade para buscar a temperatura na WeatherAPI.

Conseguindo a temperatura em Celsius, nossa entidade é então criada, e o serviço de domínio se encarrega de calcular as conversões para as demais unidades de medida (Fahrenheit e Kelvin).

Uma vez tendo as 3 temperaturas, devolvemos na API o valor das mesmas.

# Instrucoes

## Docker

Para executar o programa via docker, rode:
```sh
docker compose up -d
```

## Google Cloud Run

Acesse https://gcloud-run-strivdutlq-rj.a.run.app/{cep}, informando o CEP desejado como `Path Param`.

# Testes

Foram adicionados alguns testes unitários nos Value Objects do domínio. Para executá-los, rode:
```sh
go test -C internal/domain/vo ./...
```