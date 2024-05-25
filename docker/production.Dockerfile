FROM warehouse-base as builder

RUN go mod download

COPY . .

# Deployment container
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/warehouse/warehouse /warehouse
ENTRYPOINT ["/warehouse"]
CMD []
