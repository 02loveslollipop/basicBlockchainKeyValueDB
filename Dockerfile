FROM golang:1.23.2-bookworm
WORKDIR /src
COPY /src /src/
RUN go mod download
RUN go build -o main cmd/main/main.go
EXPOSE 8080
CMD ["./main", "/shared/host"]