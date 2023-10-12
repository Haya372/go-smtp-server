MAIL_FILE=test.txt

build:
	go build cmd/server/smtp.go

test:
	go test ./internal/*

test-coverage:
	if [ -e cover.out ]; then rm cover.out; fi
	go test ./internal/* -cover -coverprofile cover.out && go tool cover -html=cover.out

start:
	go run cmd/server/smtp.go

send-test-mail:
	curl smtp://localhost:25 --mail-from 'from@localhost' --mail-rcpt 'to@localhost' -T ${MAIL_FILE}

generate-mock-session:
	mockgen -source=internal/session/session.go -destination=./internal/mock/mock_session.go -package=mock

generate-mock-session-factory:
	mockgen -source=internal/session/factory.go -destination=./internal/mock/mock_session_factory.go -package=mock

generate-mock-command:
	mockgen -source=internal/command/handler.go -destination=./internal/mock/mock_command.go -package=mock

generate-mock-service-auth:
	mockgen -source=internal/service/auth.go -destination=./internal/mock/mock_auth_service.go -package=mock

generate-mock-all: generate-mock-session generate-mock-command generate-mock-session-factory generate-mock-service-auth
