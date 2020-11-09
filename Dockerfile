FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN go build -o /out/api

FROM alpine:latest as runner
COPY --from=builder /out/api /
ENV PORT=8000

ENTRYPOINT ["/api"]
