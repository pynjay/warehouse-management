FROM warehouse-base

RUN go install github.com/cosmtrek/air@v1.52.0
COPY . .
RUN go mod download

CMD ["air", "-c", ".air.toml"]
