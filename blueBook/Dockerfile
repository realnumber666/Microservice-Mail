FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
ENV GO111MODULE=off
ENV GOPROXY=https://goproxy.cn
CMD ["/app/main"]