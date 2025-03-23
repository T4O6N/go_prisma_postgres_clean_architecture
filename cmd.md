# INSTALL GO MOD INIT INTO PROJECT
- go mod init user-classroom-project

# LOGGING
- go get github.com/gorilla/mux

# DOTENV
- go get github.com/joho/godotenv

# SWAGGER
- go get github.com/swaggo/swag/cmd/swag
- go get github.com/swaggo/http-swagger

# PRISMA
- npm install -g prisma
- prisma init
- prisma db push
- go get github.com/prisma/prisma-client-go
- go get github.com/steebchen/prisma-client-go
- go get github.com/steebchen/prisma-client-go/engine@v0.47.0

# SWAGGER
- go get -u github.com/swaggo/swag/cmd/swag
- go install github.com/swaggo/swag/cmd/swag@latest
- swag init
- swag init -g cmd/main.go --output ./docs
- go get -u github.com/swaggo/gin-swagger
- go get -u github.com/swaggo/files

# HTTP
- go get github.com/gin-gonic/gin

- go install github.com/air-verse/air@latest

# CORS
- go get github.com/gin-contrib/cors

# REDIS
- go get github.com/redis/go-redis/v9
- go get github.com/redis/go-redis/v9/extra/redisotel

# RABBITMQ
- go get github.com/rabbitmq/amqp091-go

# library to ensure the fields are ordered correctly in the JSON response
- go get github.com/iancoleman/orderedmap

# jwt
-go get -u github.com/golang-jwt/jwt/v5

# bcrypt
- go get golang.org/x/crypto/bcrypt