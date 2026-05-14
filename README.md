# Rinha 2026 - Go

Projeto de detecção de fraude em transações utilizando KNN (K-Nearest Neighbors) e técnicas de vetorização.

## Sobre

Este projeto implementa um sistema de análise de fraude em transações financeiras usando aprendizado de máquina. O modelo utiliza vetorização de características de transações, clientes, comerciantes e terminais para classificar transações como fraudulentas ou legítimas através do algoritmo KNN.

## Estrutura do Projeto

```
.
├── cmd/
│   ├── api/              # API HTTP para scoring de fraude
│   └── preprocess/       # Ferramentas de pré-processamento
├── internal/
│   ├── knn/              # Implementação do algoritmo KNN
│   ├── server/           # Servidor HTTP
│   └── vectorizer/       # Vetorização de dados
├── resources/            # Arquivos de recursos (modelos, configurações)
├── docker-compose.yml    # Orquestração de containers
├── Dockerfile            # Imagem Docker
└── nginx.conf            # Configuração do Nginx
```

## Requisitos

- Go 1.22+
- Docker e Docker Compose (opcional)

## Instalação

### Local

```bash
# Clonar repositório
git clone <repository-url>
cd rinha-2026-go

# Baixar dependências
go mod download

# Compilar
go build -o api ./cmd/api
```

### Docker

```bash
# Build e executar com Docker Compose
docker-compose up --build
```

## Uso

### API HTTP

A aplicação expõe endpoints HTTP:

- `GET /ready` - Verifica se a API está pronta
- `POST /fraud-score` - Avalia uma transação para detectar fraude

### Exemplo de Requisição

```bash
curl -X POST http://localhost:9999/fraud-score \
  -H "Content-Type: application/json" \
  -d '{
    "id": "tx-123",
    "transaction": { ... },
    "customer": { ... },
    "merchant": { ... },
    "terminal": { ... },
    "last_transaction": { ... }
  }'
```

### Resposta

```json
{
  "approved": true,
  "fraud_score": 0.15
}
```

## Dependências Principais

- `github.com/bytedance/sonic` - Parsing e serialização JSON de alta performance
- `github.com/cloudwego/base64x` - Codificação Base64 otimizada

## Variáveis de Ambiente

- `REFS_PATH` - Caminho para arquivo binário de referências (padrão: `resources/references.bin`)

## Performance

O projeto utiliza otimizações para máxima performance:

- **Sonic** - Serialização JSON ultra-rápida com Unmarshal/Marshal otimizado
- **Base64x** - Codificação Base64 otimizada com SIMD
- KNN otimizado para buscas em alta dimensionalidade
- Pré-processamento eficiente de dados
- Layout contíguo de vetores para melhor cache locality

### Uso das Otimizações

#### JSON com Sonic

```go
import "rinha/internal/json"

// Desserializar
var req RequestBody
err := json.Unmarshal(data, &req)

// Serializar
response := ResponseBody{...}
bytes, err := json.Marshal(response)
str, err := json.MarshalToString(response)
```

#### Base64x

```go
import "rinha/internal/encoding"

// Codificar
encoded := encoding.EncodeBase64(data)
encodedURL := encoding.EncodeBase64URLSafe(data)

// Decodificar
decoded, err := encoding.DecodeBase64(encoded)
decodedURL, err := encoding.DecodeBase64URLSafe(encodedURL)
```

## Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

## Autor

Desenvolvido para Rinha de Backend 2026.
