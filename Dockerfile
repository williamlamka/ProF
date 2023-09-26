FROM golang:1.21rc2-alpine3.18
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /project-f

CMD [ "/project-f" ]