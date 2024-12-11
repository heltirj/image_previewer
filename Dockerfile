FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/image_previewer ./cmd/image_previewer

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder app/bin/image_previewer /bin/image_previewer

COPY configs/ /app/configs/

WORKDIR /app
ENTRYPOINT ["/bin/image_previewer"]