FROM golang:1.16-alpine

# ENV GO111MODULE=on
# ENV GOFLAGS=-mod=vendor

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy
# RUN go mod vendor

# COPY src/ /src/
COPY . ./
RUN ls

# RUN cd /src && \
#     go build -o /wikimedia-pageviews-api

RUN go build -o /wikimedia-pageviews-api ./src/main.go

ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

CMD [ "/wikimedia-pageviews-api" ]