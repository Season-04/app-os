FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY pb pb
COPY *.go /app/

RUN CGO_ENABLED=0 go build -o /app/clock .


FROM gcr.io/distroless/static as dist

WORKDIR /app

COPY appOS.manifest.json .

COPY --from=builder /app/clock .

EXPOSE 8080
EXPOSE 8081

CMD [ "/app/clock" ]
