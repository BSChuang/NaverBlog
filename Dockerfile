FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o naver-blog .

FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/naver-blog .

USER nonroot:nonroot

ENTRYPOINT ["/naver-blog"]