FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o notirc .

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/notirc /notirc

EXPOSE 8080

ENTRYPOINT ["/notirc"]
