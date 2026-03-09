FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache gcc musl-dev git

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags='-s -w -linkmode external -extldflags "-static"' \
    -tags 'netgo osusergo sqlite_omit_load_extension' \
    -o book2shelf ./main

FROM alpine:3.19

WORKDIR /app

# Dependency for PDF to image conversion
RUN apk add --no-cache poppler-utils ca-certificates

COPY --from=builder /build/book2shelf /app/book2shelf

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE=release
ENV TZ=UTC

ENTRYPOINT ["/app/book2shelf"]
