FROM golang:1.17.2-alpine3.14
WORKDIR /app
ADD . .
RUN apk add git 
RUN export GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/ -v ./server
FROM alpine:3.13.6
WORKDIR /app
COPY --from=0 /app/build/server .
CMD ["/app/server"]

# docker run -d ***/go-grpc-demo:v0.1 