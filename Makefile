.PHONY: dependencies check-requirements check-code check docker-up docker-down clean check-env unit-test integration-test test all run

# 기본 Go 컴파일러 설정
GO=go
GOFLAGS=-v
BINARY_NAME=portone-payment

# 소스 파일 및 디렉토리
SRC_DIRS=./...
MAIN_FILE=cmd/main.go

# 테스트 설정
TEST_TIMEOUT=10s
COVERAGE_FILE=coverage.out

# Docker 설정
DOCKER_COMPOSE=docker-compose

help: ## 사용 가능한 명령어 표시
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dependencies: ## 의존성 설치
	$(GO) mod download
	$(GO) mod tidy

check-requirements: ## 필수 요구사항 체크
	@chmod +x scripts/check_requirements.sh
	@./scripts/check_requirements.sh

check-code: ## 코드 검사
	$(GO) fmt $(SRC_DIRS)
	$(GO) vet $(SRC_DIRS)

check: check-requirements check-code ## 모든 요구사항 체크

docker-up: ## Docker 컨테이너 실행 (재빌드 포함)
	$(DOCKER_COMPOSE) up --build

docker-down: ## Docker 컨테이너 종료 및 리소스 정리
	$(DOCKER_COMPOSE) down -v --rmi local --remove-orphans

clean: ## 빌드 파일 및 캐시 제거
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)

check-env: ## 환경 변수 확인
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

unit-test: check-env ## 유닛 테스트 실행
	@sh -c "export $$(cat .env | grep -v '^#' | xargs) && $(GO) test $(GOFLAGS) -race -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_FILE) $(SRC_DIRS)"
	$(GO) tool cover -func=$(COVERAGE_FILE)

integration-test: check-env ## 통합 테스트 실행
	@sh -c "export $$(cat .env | grep -v '^#' | xargs) && $(GO) test ./test/integration/... -v -timeout $(TEST_TIMEOUT)"

test: unit-test integration-test ## 모든 테스트 실행

all: clean test run ## 전체 빌드 프로세스 실행

run: check-env ## 프로그램 실행
	@sh -c "export $$(cat .env | grep -v '^#' | xargs) && $(GO) run $(MAIN_FILE)"
