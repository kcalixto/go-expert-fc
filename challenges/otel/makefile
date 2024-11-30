.PHONY: rm build up

rm:
	docker-compose down
	docker-compose rm -f

build:
	docker-compose build

up: rm build
	docker-compose up -d
# executes some curls to check services traces
	for i in {1..1}; do \
		curl -X POST http://localhost:8080/temperature --data '{"cep":"01001000"}'; \
	done