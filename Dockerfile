FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY . .

RUN go build  -o main-app ./cmd/main.go
RUN go build  -o migrate ./cmd/migrate/main.go

CMD ["./main-app"]
