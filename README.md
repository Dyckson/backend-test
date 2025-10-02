# 🍺 Beer Temperature & Spotify Recommendation API

Uma API REST para gerenciamento de estilos de cerveja e recomendação de playlists baseada na temperatura ideal de consumo.

## � Como Executar Localmente

### Pré-requisitos

- Docker e Docker Compose instalados

### Setup Rápido

```bash
# 1. Clone o repositório
git clone https://github.com/Dyckson/backend-test
cd backend-test

# 2. Configure as variáveis de ambiente
cp .env

# 3. Execute com Docker
docker-compose up --build

# ✅ API disponível em: http://localhost:1112
```

## 📚 Como Usar a API

### Base URL

```
http://localhost:1112/api
```

### 🍺 Estilos de Cerveja (CRUD)

#### Listar todos os estilos

```bash
curl -X GET http://localhost:1112/api/beer-styles/list
```

#### Criar novo estilo

```bash
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "IPA",
    "temp_min": -6.0,
    "temp_max": 7.0
  }'
```

#### Atualizar estilo

```bash
curl -X PUT http://localhost:1112/api/beer-styles/edit/{uuid} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Double IPA",
    "temp_min": -7.0,
    "temp_max": 8.0
  }'
```

#### Deletar estilo

```bash
curl -X DELETE http://localhost:1112/api/beer-styles/{uuid}
```

### 🎵 Recomendação de Playlist

#### Obter recomendação baseada na temperatura

```bash
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": -7.0}'
```

**Resposta:**

```json
{
  "beerStyle": "IPA",
  "playlist": {
    "name": "Rock Playlist for IPA",
    "tracks": [
      {
        "name": "Bohemian Rhapsody",
        "artist": "Queen",
        "link": "https://open.spotify.com/track/4u7EnebtmKWzUH433cf5Qv"
      }
    ]
  }
}
```

## 🧪 Executar Testes

```bash
# Todos os testes
go test ./... -v

# Apenas unit tests
go test ./internal/http/controller/ -v

# Apenas integration tests
go test ./tests/integration/ -v
```

## � Tecnologias

- **Go 1.24.5** com Gin framework
- **PostgreSQL**
- **Spotify Web API**
- **Docker**
- **Clean Architecture** com testes híbridos

## 📋 Status Codes

| Código | Descrição               |
| ------- | ------------------------- |
| `200` | Sucesso                   |
| `201` | Criado                    |
| `400` | Dados inválidos          |
| `404` | Não encontrado           |
| `409` | Conflito (nome duplicado) |
| `500` | Erro interno              |
| `503` | Spotify indisponível     |

---

**Para documentação completa:** consulte os arquivos `DEVELOPMENT.md`, `API.md` e `FEATURES.md`
