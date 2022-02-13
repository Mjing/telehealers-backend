server.out: src/*/*.go swagger.obj
	go build -o $@ cmd/telehealers-backend-server/main.go

swagger.obj: swagger/swagger.yml
	swagger generate server -A telehealers-backend -f ./swagger/swagger.yml
	echo "" > $@

clean:
	rm *.out *.obj