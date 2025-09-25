FROM golang:1.22
WORKDIR /app
COPY ../scaffold .
COPY ../.env.docker .env
RUN chmod +x scaffold
CMD ["./scaffold"]

#$env:GOOS="linux"
#$env:GOARCH="amd64"
#go build -o scaffold .
