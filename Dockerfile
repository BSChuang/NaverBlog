FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o quiz-bot .

FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/quiz-bot .

USER nonroot:nonroot

ENTRYPOINT ["/quiz-bot"]