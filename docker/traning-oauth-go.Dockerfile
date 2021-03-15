FROM golang
LABEL maintainer="Hebert <hebert.t.santos@gmail.com>"
WORKDIR /app/src/traning-oauth-go
ENV GOPATH=/app
COPY . /app/src/traning-oauth-go
RUN go build main.go
ENTRYPOINT ["./main"]
CMD ["-conf", "config.yaml"]
EXPOSE 8082