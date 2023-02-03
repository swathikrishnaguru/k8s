FROM golang:alpine

WORKDIR /app
COPY . .

RUN go build -o scaler

CMD ["./scaler"]
