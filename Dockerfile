FROM golang:1.20-alpine

ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy

COPY . ./

RUN go build -o /wikimedia-pageviews-api ./src/main.go

ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

CMD [ "/wikimedia-pageviews-api" ]