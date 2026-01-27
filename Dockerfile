FROM golang:1.25.5

WORKDIR /app

# Сначала зависимости — для кеша
COPY go.mod go.sum ./
RUN go mod download

# Потом весь код
COPY . .

RUN go build -o app

EXPOSE 8080

CMD ["./app"]
