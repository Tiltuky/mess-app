# Elasticsearch — это система полнотекстового поиска, написанная на Java. 
### Кроме того, Elasticsearch — это нереляционное хранилище документов в формате JSON, которое выпущено как проект с открытым исходным кодом в соответствии с условиями лицензии Apache. 
## [Elastic guide](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html)
# Создание пользователя: 
нужно совершить вставку записи в postgres и elastic
postgres: операция INSERT
elastic: client.Index

#### Elastic: Indexing documents 
```
document := struct {
    Name string `json:"name"`
}{
    "go-elasticsearch",
}
data, _ := json.Marshal(document)
client.Index("my_index", bytes.NewReader(data))
```
# Обновление пользователя: 
postgres: операция UPDATE
elastic: client.Update
#### Elastic: Update
```
client.Update("my_index", "id", strings.NewReader(`{doc: { language: "Go" }}`))
```
# Удаление пользователя: 
postgres: операция DELETE
elastic: client.Delete
#### Elastic: Delete
```
client.Delete("my_index", "id")
```
# Поиск пользователя: 
postgres: операция DELETE
elastic: client.Search  

> Elastic используется как поисковой движок, поэтому поисковой запрос выполняется в нем, и по полученным идентификаторам данных запрашиваем из Postgres. 
#### Elastic: Search
```
query := `{ "query": { "match_all": {} } }`
client.Search(
    client.Search.WithIndex("my_index"),
    client.Search.WithBody(strings.NewReader(query)),
)
```

## Install
``` 
go get github.com/elastic/go-elasticsearch/v8@latest 
```
## Connect
```
client, err := elasticsearch.NewClient(elasticsearch.Config{
    CloudID: "<CloudID>",
    APIKey: "<ApiKey>",
})
```

## Create index
```
client.Indices.Create("my_index")
```
## Get
```
client.Get("my_index", "id")
```

# Elasticsearch и Postgres - почему вместе?
### Elasticsearch: 
Идеально подходит для полнотекстового поиска, аналитики и работы с неструктурированными данными.
### PostgreSQL: 
Превосходит в обработке транзакций, сложных запросов к структурированным данным и поддержании целостности данных.
### Комбинированные рабочие процессы: 
Многие приложения требуют как быстрого поиска по большому объему текстовой информации, так и надежного хранения структурированных данных. Комбинируя Elasticsearch и PostgreSQL, можно эффективно решать такие задачи.



# Docker-compose
```
services:
  elasticsearch:
    image: 'docker.elastic.co/elasticsearch/elasticsearch:8.1-SNAPSHOT'
    build:
      context: ./docker/es
      dockerfile: Dockerfile
      args:
        DISCOVERY_TYPE: single-node
    ports:
      - '9200:9200'
    volumes:
      - elasticdata:/usr/share/elasticsearch/data
```
# .env
#### environment variable used by elasticsearch.NewDefaultClient()
```
ELASTICSEARCH_URL=http://elasticsearch:9200
```


