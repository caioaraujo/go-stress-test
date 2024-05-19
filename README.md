# GO Stress Test
Projeto CLI para testes de carga em um servidor WEB.

## Execução local
`go run main.go --url={URL} --requests={REQUESTS} --concurrency={CONCURRENCY}`

Onde:

- requests: Total de requisições permitidas;
- concurrency: Total de concorrências onde as requisições serão distribuídas.

### Docker
Alternativamente é possível executar via Docker:

`docker run stress-test --url={URL} --requests={REQUESTS} --concurrency={CONCURRENCY}`
