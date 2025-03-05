# Serviço de Clima por CEP

Este serviço fornece informações meteorológicas com base em CEPs brasileiros. Ele retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## URL do Serviço

O serviço está disponível em:
```
https://go-cep-clima-727490326131.us-central1.run.app/
```

Exemplo de uso:
```
GET https://go-cep-clima-727490326131.us-central1.run.app/weather?cep=01001000
```

## Requisitos

- Go 1.23 ou superior
- Docker e Docker Compose
- Chave de API do WeatherAPI (cadastre-se em https://www.weatherapi.com/)

## Executando Localmente

1. Configure sua chave da WeatherAPI como variável de ambiente:
```bash
export WEATHER_API_KEY=sua_chave_api_aqui
```

2. Execute com Docker Compose:
```bash
docker-compose up --build
```

O serviço estará disponível em http://localhost:8080

## Uso da API

### Consultar Clima por CEP

```
GET /weather?cep=12345678
```

#### Resposta de Sucesso (200 OK)
```json
{
    "temp_C": 25.0,
    "temp_F": 77.0,
    "temp_K": 298.15
}
```

#### Respostas de Erro

- CEP com formato inválido (422 Unprocessable Entity)
```json
{
    "message": "invalid zipcode"
}
```

- CEP não encontrado (404 Not Found)
```json
{
    "message": "can not find zipcode"
}
```

## Implantação

O serviço foi projetado para ser implantado no Google Cloud Run. Siga estas etapas:

1. Construa e envie a imagem Docker:
```bash
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/cep-weather
```

2. Implante no Cloud Run:
```bash
gcloud run deploy cep-weather \
  --image gcr.io/SEU_PROJECT_ID/cep-weather \
  --platform managed \
  --set-env-vars WEATHER_API_KEY=sua_chave_api_aqui
```

## Funcionalidades

- Validação de CEP (8 dígitos)
- Integração com a API ViaCEP para busca de localidades
- Integração com a WeatherAPI para dados meteorológicos
- Conversão automática de temperaturas entre Celsius, Fahrenheit e Kelvin
- Tratamento adequado de erros
- Containerização com Docker
- Pronto para deploy no Google Cloud Run 