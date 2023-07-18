FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY pb pb
COPY *.go /app/

RUN go build -o /app/clock .

EXPOSE 8080
EXPOSE 8081

CMD [ "/app/clock" ]
