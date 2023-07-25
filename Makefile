user-service.docker.start:
	docker compose --file ./users-service/docker-compose.yml  up -d

run.user-service:  user-service.docker.start 
		go run ./users-service/cmd/main.go
