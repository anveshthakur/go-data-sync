BINARY_NAME=syncData

build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/web
	@echo "Built!"

run: build
	@echo "Starting..."
	@env ./${BINARY_NAME} &
	@echo "Started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

start: run

start_ui:
	@echo "Starting UI..."
	@cd ui && npm run dev
	@echo "Started!"

stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"
