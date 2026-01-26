FROM golang:1.25.5

WORKDIR /app

# Сначала зависимости — для кеша
COPY go.mod go.sum ./
RUN go mod download

# Потом весь код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080

CMD ["./app"]
