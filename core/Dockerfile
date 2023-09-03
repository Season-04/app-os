FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY internal internal
COPY middleware middleware
COPY types types
COPY manifest manifest
COPY *.go /app/

RUN CGO_ENABLED=0 go build -o /app/core .


FROM gcr.io/distroless/static as dist

WORKDIR /app

COPY appOS.manifest.json .

COPY --from=builder /app/core .

EXPOSE 8080

CMD [ "/app/core" ]
