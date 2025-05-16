FROM golang:1.24

WORKDIR /app

RUN apt install -y make
RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
