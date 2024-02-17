FROM golang:1.22-alpine

WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -C main -o /rinha-backend-2024-q1

EXPOSE 9091

CMD ["/rinha-backend-2024-q1"]
