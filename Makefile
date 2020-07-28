build:
	@echo Building binary...
	@CGO_ENABLED=0 GOOS=linux go build -o ./dist/datahow ./cmd/*.go

clean:
	@echo Removing binary...
	@rm ./dist/datahow
