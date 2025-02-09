FROM golang:1.22.10 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go


FROM debian AS runner

RUN apt update && apt install -y iproute2 iputils-ping procps

COPY --from=builder main /bin/main

ENTRYPOINT ["/bin/main"]
