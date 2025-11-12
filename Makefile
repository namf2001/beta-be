.PHONY: build-local-go-image api-setup api-run api-pg-migrate-up api-pg-migrate-down api-gen-models api-go-generate api-gen-mocks pg

DOCKER_BIN := docker
PROJECT_NAME := beta-be
DOCKER_COMPOSE_BIN := docker-compose
DOCKER_BIN := docker

COMPOSE := PROJECT_NAME=${PROJECT_NAME} ${DOCKER_COMPOSE_BIN} -f build/docker-compose.base.yaml -f build/docker-compose.local.yaml
API_COMPOSE = ${COMPOSE} run --name ${PROJECT_NAME}-api-local --rm --service-ports -w /api api

build-local-go-image:
	${DOCKER_BIN} build -f build/local.go.Dockerfile -t ${PROJECT_NAME}-api-local:latest .
	-${DOCKER_BIN} images -q -f "dangling=true" | xargs ${DOCKER_BIN} rmi -f
api-setup: pg api-pg-migrate-up
	sleep 5
	${DOCKER_BIN} image inspect ${PROJECT_NAME}-go-local:latest >/dev/null 2>&1 || make build-local-go-image
api-run:
	${API_COMPOSE} sh -c "go run -mod=vendor cmd/server/main.go server -c configs/.env"
api-pg-migrate-up:
	${COMPOSE} run --rm pg-migrate sh -c './migrate -path /api-migrations -database $$PG_URL up'
api-pg-migrate-down:
	${COMPOSE} run --rm pg-migrate sh -c './migrate -path /api-migrations -database $$PG_URL drop'
api-gen-models:
	${API_COMPOSE} sh -c 'cd ./internal/repository && go run ariga.io/entimport/cmd/entimport -dsn "postgres://${PROJECT_NAME}:@pg:5432/${PROJECT_NAME}?sslmode=disable" && go run entgo.io/ent/cmd/ent generate --feature sql/execquery ./ent/schema'
api-go-generate:
	${API_COMPOSE} sh -c "go generate ./..."
api-gen-mocks:
	${COMPOSE} run --name ${PROJECT_NAME}-mockery-local --rm -w /api --entrypoint '' mockery /bin/sh -c "\
		mockery --dir internal/controller --all --recursive --inpackage && \
		mockery --dir internal/repository --all --recursive --inpackage"
pg:
	${COMPOSE} up -d pg
