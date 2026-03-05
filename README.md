# weather-go

Aplicação em Go para consulta de dados meteorológicos utilizando provedores externos de clima.

O projeto permite consultar condições atuais e previsão do tempo utilizando:

- WeatherAPI (requer token)
- Open-Meteo (não requer token)

A aplicação é executada via Docker Compose, simplificando o setup e a execução.

---

## Sobre o projeto

`weather-go` é um serviço distribuído que consome APIs públicas de meteorologia e processa dados climáticos de forma estruturada.

Atualmente suporta dois provedores:

- WeatherAPI (https://www.weatherapi.com/)
- Open-Meteo (https://open-meteo.com/)

O provedor é definido via configuração, permitindo alternância sem necessidade de alteração de código. A arquitetura foi projetada para ser extensível, possibilitando a implementação e integração de novos provedores de forma desacoplada.

### Arquitetura

O sistema é composto por três serviços principais:

- **receiver**  
  Responsável por receber os dados coletados pelos demais serviços e persistir as informações em um banco de dados temporal.

- **collector**  
  Coleta dados meteorológicos históricos e atuais a partir do provedor configurado e os envia ao `receiver`.

- **forecast**  
  Coleta dados de previsão meteorológica e os encaminha ao `receiver`.

Todos os serviços são executados via Docker Compose.

### Persistência

Os dados são armazenados em um banco de dados temporal utilizando [InfluxDB](https://www.influxdata.com/), permitindo consultas eficientes baseadas em séries temporais.

---

## Executando localmente

### Pré-requisitos

Para executar o projeto é necessário:

- Docker
- Docker Compose
- Token da WeatherAPI (apenas se utilizar o provedor `weatherapi`)

Não é necessário ter o Go instalado localmente, pois todos os serviços são executados via containers.

### 1. Clonar o repositório

```bash
git clone https://github.com/cardosodg/weather-go.git
cd weather-go
```

### 2. Configuração do projeto
Crie o arquivo `.env` a partir do modelo:
```bash
cp .env.example .env
```

Edite o arquivo .env com os valores adequados:
```env
# Token da API WeahterAPI
WEATHER_API_KEY=************************

# Configurações do banco de dados
INFLUXDB_MODE=setup
INFLUXDB_USERNAME=usuario
INFLUXDB_PASSWORD=sua_senha
INFLUXDB_ORG=monitoring
INFLUXDB_BUCKET=weather
INFLUXDB_ADMIN_TOKEN=seu_influxdb_token
INFLUXDB_URL="http://influxdb:8086"
```
### Lista de variáveis de ambiente

| Variável              | Obrigatória | Descrição |
|-----------------------|------------|------------|
| WEATHER_API_KEY       | Apenas para `weatherapi` | Token de autenticação da WeatherAPI |
| INFLUXDB_MODE         | Sim        | Modo de execução do InfluxDB (ex: setup, init, etc.) |
| INFLUXDB_USERNAME     | Sim        | Usuário administrativo do InfluxDB |
| INFLUXDB_PASSWORD     | Sim        | Senha do usuário administrativo |
| INFLUXDB_ORG          | Sim        | Organização utilizada no InfluxDB |
| INFLUXDB_BUCKET       | Sim        | Bucket onde os dados meteorológicos serão armazenados |
| INFLUXDB_ADMIN_TOKEN  | Sim        | Token administrativo do InfluxDB |
| INFLUXDB_URL          | Sim        | URL de conexão com o InfluxDB (ex: http://influxdb:8086) |

### Lista de localidades

A lista de localidades define quais coordenadas serão monitoradas pelos serviços `collector` e `forecast`.

Crie o arquivo `locations.json` a partir do modelo:

```bash
cp locations.example.json locations.json
```

A estrutura esperada de locations.json é:
```json
[
  {
    "latitude": "-1.123456",
    "longitude": "-2.123456",
    "label": "localA"
  },
  {
    "latitude": "-2.345678",
    "longitude": "-3.345678",
    "label": "localB"
  }
]
```
O campo `label` é utilizado pelo `receiver` como tag nas medições armazenadas no InfluxDB, permitindo identificar cada localidade nas consultas.

E `locations.json` pode mudar de acordo com a necessidade de se adicionar ou remover localidades de monitoramento.

Para definir qual provedor será utilizado pelos serviços `collector` e `forecast`, edite os seguintes arquivos:
```
internal/collector/config/config.go
internal/forecast/config/config.go
```

Em ambos, altere a variável `Provider` conforme desejado:

```go
// possible values: openmeteo or weatherapi
Provider = "weatherapi"
```

Valores possíveis:
- weatherapi
- openmeteo

A alteração deve ser feita nos dois serviços para garantir consistência no provedor utilizado.

---

### 3. Executar os containers
Com as configurações realizadas, inicie os serviços principais com:
```bash
docker compose up -d collector receiver influxdb
```

O serviço `forecast` pode ser executado de forma agendada ou manual.

### Execução via cron
Sugestão de agendamento no `crontab`
```bash
0 9,21 * * * cd path/to/weather-go && docker compose run --rm forecast >> /var/log/forecast.log 2>&1
```
O exemplo acima executa o forecast duas vezes ao dia, às 9h e 21h.

### Execução manual

```bash
docker compose run --rm forecast
```

### Verificar logs
Para acompanhar os logs dos serviços em execução:
```bash
docker compose logs -f collector receiver
```

---

## Visualização dos dados

Para criação de dashboards e visualização das métricas coletadas, pode-se utilizar o [Grafana](https://grafana.com/).

O Grafana pode ser conectado diretamente ao InfluxDB. Basta adicionar uma nova *Data Source* do tipo **InfluxDB** e configurar a conexão apontando para o container do InfluxDB.

Para InfluxDB 2.x, será necessário informar:

- URL (ex: `http://influxdb:8086`)
- Organization (`INFLUXDB_ORG`)
- Bucket (`INFLUXDB_BUCKET`)
- Token (`INFLUXDB_ADMIN_TOKEN`)

### Exemplo de serviço Grafana no Docker Compose

Você pode adicionar o seguinte serviço ao seu `docker-compose.yml`:

```yaml
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    restart: unless-stopped
```

Após subir o container:

```bash
docker compose up -d grafana
```

O Grafana estará disponível em:

```
http://localhost:3000
```

Credenciais padrão:

- **Usuário:** admin  
- **Senha:** admin  

Lembre-se de trocar as credenciais padrão do Grafana.
