FROM golang:1.21.6

WORKDIR /letsgo

RUN git clone https://github.com/PranaySiddharthBhange/letsgo.git .

CMD ["go", "run", "main.go"]
