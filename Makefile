build: 
	nmake user
	go build main.go
	
run:
	nmake run-user
	nmake run-auth

run-user:
	go run main.go

run-auth:
	cd ../auth-service/
	go run main.go

user:
	protoc -I. --go_out=plugins=grpc:./services/user --grpc-gateway_out=logtostderr=true:./services/user --swagger_out=logtostderr=true:./services/user --proto_path proto user.proto 
	mv ./services/user/egnite.app/microservices/user/proto/user/* ./services/user/
	rm -r ./services/user/egnite.app
	cp ./services/user/*.json ./www/