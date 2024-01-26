FROM golang:1.21.6

WORKDIR /letsgo

COPY . .

CMD ["go", "run", "main.go"]
