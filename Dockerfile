FROM golang:alpine as builder
# Install git
RUN apk update && apk add git
RUN mkdir /go/src/kernel-concierge/
WORKDIR /go/src/kernel-concierge/
ADD . /go/src/kernel-concierge/
# get dependencies
RUN go get -d -v
# compile code
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /go/src/kernel-concierge/ /app/
WORKDIR /app
CMD ["./main"]