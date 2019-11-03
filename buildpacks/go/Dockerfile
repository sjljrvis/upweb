FROM golang:1.12

RUN go get github.com/Masterminds/glide

WORKDIR /go/src/

ADD glide.yaml glide.yaml
ADD glide.lock glide.lock

RUN glide install

ADD . .

CMD ["go", "run", "main.go"]