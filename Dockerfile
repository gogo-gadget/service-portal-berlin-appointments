FROM golang:1.18.2-alpine3.15 as builder

WORKDIR /app

COPY ./ /app/

RUN go build -o /app/dist/service-portal-berlin-appointments

FROM alpine:3.15.4

WORKDIR /app

COPY --from=builder /app/dist/service-portal-berlin-appointments /app/service-portal-berlin-appointments
COPY --from=builder /app/config/.config.yaml /app/config/.config.yaml

CMD ["/app/service-portal-berlin-appointments"]
