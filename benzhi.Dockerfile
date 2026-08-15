FROM golang:1.26.2

WORKDIR /workspace
ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["go", "run", "./cmd/manifestctl"]
