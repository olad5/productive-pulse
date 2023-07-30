
user-service.docker.start:
	docker compose --file ./users-service/docker-compose.yml  up -d

users-service.docker.test.start:
	docker compose --file ./users-service/tests/config/docker-compose.yml  up -d

test.users-service: users-service.docker.test.start
		go test -count=1 -tags integration  ./users-service/tests/integration/

test.users-service.verbose: users-service.docker.test.start
		go test -count=1 -tags integration -v ./users-service/tests/integration/

run.user-service:  user-service.docker.start 
		go run ./users-service/cmd/main.go




todo-service.docker.start:
	docker compose --file ./todo-service/docker-compose.yml  up -d

todos-service.docker.test.start:
	docker compose --file ./todo-service/tests/config/docker-compose.yml  up -d

test.todos-service: todos-service.docker.test.start
		go test -count=1 -tags integration  ./todo-service/tests/integration/

test.todos-service.verbose: todos-service.docker.test.start
		go test -count=1 -tags integration -v ./todo-service/tests/integration/

run.todo-service:  todo-service.docker.start 
		go run ./todo-service/cmd/main.go

