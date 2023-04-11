#build stage
FROM golang:alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download 

FROM golang:alpine AS builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/app -v ./cmd/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg
COPY --from=builder /go/src/app/templates ./templates 
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=auth Version=0.0.1
EXPOSE 1323
