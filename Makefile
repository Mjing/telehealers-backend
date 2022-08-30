server.out: src/*/*.go src/*/*/*.go swagger.obj restapi/configure_telehealers_backend.go
	go build -o $@ cmd/telehealers-backend-server/main.go

swagger.obj: swagger/*.yml
	swagger flatten swagger/swagger.yml --format=yaml > flattened_swagger.yml
	swagger generate server -A telehealers-backend -P models.Principal -f ./flattened_swagger.yml
	echo "" > $@

#NOTE: Load env before running test
run_tests:
	go clean -testcache
	go test telehealers.in/router/src/swagger_service_handler/db_apis
clean:
	rm -rf *.out *.obj restapi/operations flattened_swagger.yml
	rm -rf models restapi/doc.go restapi/embedded_spec.go restapi/server.go
