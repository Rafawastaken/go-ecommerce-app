server:
	cross-env APP_ENV=dev nodemon --ext go,json \
		--watch ./cmd --watch ./internal --watch ./config \
		--exec go run cmd/main.go
