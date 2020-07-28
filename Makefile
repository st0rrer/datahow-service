build:
	@echo Building binary...
	@go build -o ./dist/datahow ./cmd/*.go

clean:
	@echo Removing binary...
	@rm ./dist/datahow
