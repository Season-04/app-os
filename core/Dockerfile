FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY pb pb
COPY users users
COPY *.go /app/

RUN go build -o /app/core .

EXPOSE 8080

CMD [ "/app/core" ]
