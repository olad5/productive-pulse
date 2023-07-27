user-service.docker.start:
	docker compose --file ./users-service/docker-compose.yml  up -d

run.user-service:  user-service.docker.start 
		go run ./users-service/cmd/main.go


users-service.docker.test.start:
	docker compose --file ./users-service/tests/config/docker-compose.yml  up -d

test.users-service: users-service.docker.test.start
		go test -count=1 -tags integration  ./users-service/tests/integration/

test.users-service.verbose: users-service.docker.test.start
		go test -count=1 -tags integration -v ./users-service/tests/integration/
