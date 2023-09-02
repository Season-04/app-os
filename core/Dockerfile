FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY internal internal
COPY middleware middleware
COPY types types
COPY manifest manifest
COPY *.go /app/

RUN go build -o /app/core .

COPY appOS.manifest.json .

EXPOSE 8080

CMD [ "/app/core" ]
