FROM golang:alpine3.22

RUN go version

ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY ./ ./

RUN go mod download


RUN go build -ldflags="-w -s" -o url-shortener ./cmd/main.go

RUN chmod +x ./url-shortener

CMD ["./url-shortener"]