FROM golang:1.18.4-alpine

WORKDIR /go/src/github.com/meidonohitsuji/gitlab-telegram-notification-go

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /go/src/github.com/meidonohitsuji/gitlab-telegram-notification-go

RUN go build -o /gitlab-telegram-notification-go

CMD [ "/gitlab-telegram-notification-go" ]