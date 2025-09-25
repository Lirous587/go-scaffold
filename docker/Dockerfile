FROM golang:1.22
WORKDIR /app
COPY ../main .
COPY ../.env.docker .env
RUN chmod +x main
CMD ["./main"]

#$env:GOOS="linux"
#$env:GOARCH="amd64"
#go build -o main .
