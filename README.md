# go-load-balancer

Балансировщик нагрузки. Первая и вторая части [тестового задания](https://github.com/Go-Cloud-Camp/test-assignment/tree/25334d85c4e90cccc9cf0f93bdd275738295ad13).

## Демо через docker compose

Предварительно нужно освободить порты `8080` и `9090` на хосте, либо поменять их в `docker-compose.yml` и `config/example_docker/config.yml`.

Скачать исходный код и поднять демонстрационный стенд (соберется образ балансировщика, скачаются `nginx:stable-alpine3.21` и `postgres:16-alpine3.21`):

```shell
git clone https://github.com/t0ggit/go-load-balancer.git
cd go-load-balancer
```

```shell
docker compose build --no-cache
```

```shell
docker compose up -d
```

Смотреть логи балансировщика в реальном времени:

```shell
docker compose logs load-balancer --tail=1000 -f
```

Теперь можно обращаться к балансировщику по адресу `localhost:8080`.

Например, использовать Apache Bench:

```shell
sudo apt install -y apache2-utils

ab -n 5000 -c 1000 http://localhost:8080/
```

Выключить стенд:

```shell
docker compose down
```

## Использование API RateLimiter'а

Получить настройки бакета для клиента с таким ключом (это может быть как `RemoteAddr`, так и API-ключ из заголовка `Authorization`)

```shell
curl -X POST http://localhost:9090/get \
  -H "Content-Type: application/json" \
  -d '{"key": "test_api_key1"}'
```

Установить новые настройки бакета для клиента:

```shell
curl -X POST http://localhost:9090/set \
-H "Content-Type: application/json" \
-d '{
    "key": "test_api_key1",
    "bucket_capacity": 100,
    "refill": 50,
    "refill_interval": "5s"
}'
```

Обратиться к балансировщику с заголовком `Authorization`:

```shell
curl -H "Authorization: test_api_key1" http://localhost:8080/
```

```shell
ab -n 120 -c 1 -H "Authorization: test_api_key1" http://localhost:8080/ | grep requests:
```

## UML диаграмма

https://www.dumels.com/

![image](https://github.com/user-attachments/assets/aed1661e-8062-42bc-86b2-1271d8755ae1)
