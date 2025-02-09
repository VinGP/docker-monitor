FROM golang:1.22.10 AS builder

WORKDIR /build

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN swag init --parseDependency --parseInternal -g ./internal/http/handlers/handlers.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

FROM alpine AS runner

RUN apk add --no-cache curl

COPY --from=builder main /bin/main

ENTRYPOINT ["/bin/main"]
