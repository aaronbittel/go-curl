BUILD_DIR = build
EXE_NAME = go-curl

.PHONY = build

build:
	go build -o $(BUILD_DIR)/$(EXE_NAME)

run: build
	$(BUILD_DIR)/$(EXE_NAME)
