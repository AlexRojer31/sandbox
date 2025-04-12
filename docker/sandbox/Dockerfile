FROM golang:1.23-alpine3.21

WORKDIR /sandbox
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY configs ./configs/
COPY internal ./internal/

EXPOSE 8082

ENTRYPOINT ["go", "run" , "cmd/sandbox/main.go"]
CMD ["--help"]