# Путь к проекту и бинарникам
PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin

# Путь и версия генератора моков
MOQ = $(PROJECT_BIN)/moq
MOQ_VERSION = v0.3.1

# === Mocks generator ===

.PHONY: .install-moq
.install-moq:
	@echo "Installing moq..."
	@mkdir -p $(PROJECT_BIN)
	[ -f $(MOQ) ] || GOBIN=$(PROJECT_BIN) go install github.com/matryer/moq@$(MOQ_VERSION)

# Генерация моков
.PHONY: generate-mocks
generate-mocks: .install-moq
	@echo "Generating mocks..."
	$(MOQ) -out ./internal/mocks/mock_file.go ./internal/interfaces InterfaceName

# === Тесты ===

# Запуск всех тестов
.PHONY: test
test:
	go test ./...

# Быстрый запуск тестов
.PHONY: test-fast
test-fast:
	go test -v -short ./...
