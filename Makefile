.PHONY: all build clean deploy

LAMBDA_FUNCTIONS := create get update delete sqs

build:
	@mkdir -p build
	@for func in $(LAMBDA_FUNCTIONS); do \
		echo "Building $$func..."; \
		mkdir -p build/$$func; \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$$func/bootstrap ./cmd/$$func; \
		(cd build/$$func && zip -q ../$$func.zip bootstrap); \
	done
	@echo "Build complete."

clean:
	rm -rf build
	@echo "Cleaned up."

deploy: build
	serverless deploy

deploy-%:
	mkdir -p build/$*
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$*/bootstrap ./cmd/$*
	(cd build/$* && zip -q ../$*.zip bootstrap)
	serverless deploy function --function $*
