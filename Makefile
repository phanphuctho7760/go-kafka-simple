init:
	cd producer && go mod vendor
	cd consumer && go mod vendor
build:
	docker compose up -d --build
start:
	docker compose start
restart:
	docker compose restart
stop:
	docker compose stop
info:
	docker compose ps
clear:
	docker compose down -v
send:
	docker compose exec go-kafka-simple-producer /bin/project
receive:
	docker compose exec go-kafka-simple-consumer /bin/project
exec-producer:
	docker compose exec go-kafka-simple-producer sh
exec-consumer:
	docker compose exec go-kafka-simple-consumer sh
hello:
	echo "Hello '$(name)'. This is example simple Kafka with Golang"
