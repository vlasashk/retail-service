# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до нашего прото
LOMS_PROTO_PATH = ../proto_contracts/loms/v1
PROTO_GEN_OUT = pkg/api/loms/v1

.PHONY: build
build: .protoc-generate-server
	go build -o $(LOCAL_BIN)/loms ./cmd/loms/main.go

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml ./...

# Сперва сгенерировать proto контракт .protoc-generate-server
.PHONY: unit-test
unit-test: cover-folder
	go test -race -coverprofile ./coverage/cover.out ./internal/... && \
    go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
    open ./coverage/cover.html && \
    rm ./coverage/cover.out

.PHONY: integration-test
#https://github.com/golang/go/issues/65653
integration-test: export GOEXPERIMENT=nocoverageredesign
integration-test: .protoc-generate-server cover-folder
	go test -race -tags=integration -coverprofile ./coverage/cover.out -coverpkg ./internal/... ./... && \
	go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
	open ./coverage/cover.html && \
	rm ./coverage/cover.out

.PHONY: cover-folder
cover-folder:
	mkdir -p ./coverage

# Удаляем директорию vendor-proto
.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

.PHONY: sqlc-gen
sqlc-gen:
	sqlc generate -f sqlc.yaml

# Устанавливаем необходимые бинарники для компиляции
# Для начала установить protoc через brew install protobuf или по-другому под вашу ОС
# При выполнении команды может быть ошибка "GOPROXY list is not the empty string, but contains no entries"
# В этом случае нужно выполнить export GOPROXY=direct
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.21.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

# Вендоринг внешних proto файлов
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

# Генерация кода через protoc
# Важно указать наши прото файлы и proto в зависимостях через инклюды "-I"
# Также указываем плагины и директории для результата работы плагинов
.PHONY: .protoc-generate-server
.protoc-generate-server: .bin-deps .vendor-proto
	mkdir -p api/openapiv2
	mkdir -p ${PROTO_GEN_OUT}
	protoc \
	-I ${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out ${PROTO_GEN_OUT} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out ${PROTO_GEN_OUT} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/api/loms/v1" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out ${PROTO_GEN_OUT} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out api/openapiv2 \
	--openapiv2_opt logtostderr=true \
	$(LOMS_PROTO_PATH)/loms.proto
	go mod tidy

.PHONY: .serve-swagger
.serve-swagger:
	bin/swagger serve api/openapiv2/notes.swagger.json
