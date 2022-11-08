FROM golang:latest

WORKDIR /banannza

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o main .

CMD [ "/bananza" ]