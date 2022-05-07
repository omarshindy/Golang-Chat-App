FROM golang:1.18.1 as builder
# Define build env
ENV GOOS linux

ENV CGO_ENABLED 0

WORKDIR /app

RUN ls

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

FROM alpine:3.14 as production

RUN apk add --no-cache ca-certificates

COPY --from=builder app .

HEALTHCHECK  --retries=5  --interval=35s --timeout=4s CMD curl -f https://localhost/ || exit 1

EXPOSE 8080

CMD ./app