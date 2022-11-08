FROM golang:latest

COPY . /bananza

WORKDIR /bananza

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD [ "go", "run", "main.go" ]
