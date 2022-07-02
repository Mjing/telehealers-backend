server.out: src/*/*.go src/*/*/*.go swagger.obj restapi/configure_telehealers_backend.go
	go build -o $@ cmd/telehealers-backend-server/main.go

swagger.obj: swagger/swagger.yml
	swagger generate server -A telehealers-backend -f ./swagger/swagger.yml
	echo "" > $@

clean:
	rm -rf *.out *.obj restapi/operations
