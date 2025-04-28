FROM docker.io/golang:latest as builder
WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

FROM docker.io/alpine:3

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/server"]
