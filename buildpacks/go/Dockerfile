FROM golang:1.12

RUN go get github.com/Masterminds/glide

WORKDIR /go/src/dnow_replace_me

ADD glide.yaml glide.yaml
ADD glide.lock glide.lock

RUN glide install

ADD . .

CMD ["go", "run", "main.go"]