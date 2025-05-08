# go-load-balancer

Балансировщик нагрузки. Первая часть [тестового задания](https://github.com/Go-Cloud-Camp/test-assignment?tab=readme-ov-file#часть-1-балансировщик-нагрузки).

## Демо через docker compose

Предварительно нужно освободить порт `8080` на хосте, либо поменять его в `docker-compose.yml` и `config/example_docker/config.yml`.

Скачать исходный код и поднять демонстрационный стенд (соберется образ `go-load-balancer-load-balancer:latest` (до 20MB), скачается `nginx:stable-alpine3.21` (до 20MB)):

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

## Без докера

Можно собрать на хосте бинарник с детектированием состояний гонки (понадобится `gcc`)

```shell
CGO_ENABLED=1 CONFIG_PATH=config/example_config.yml go run -race ./cmd/lb/main.go | grep RACE --context=10
```

Попробовать нагрузить балансировщик (5000 запросов, 1000 соединений):

```shell
ab -n 5000 -c 1000 http://localhost:8080/
```

У меня `grep` не показал состояний гонки.

## UML диаграмма

https://www.dumels.com/

![image](https://github.com/user-attachments/assets/aed1661e-8062-42bc-86b2-1271d8755ae1)
