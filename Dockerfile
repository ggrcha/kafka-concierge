FROM golang:1.11.5 as builder

RUN apt install git
#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
ADD /main.go $GOPATH/src/kernel-concierge/main.go
RUN go get -v kernel-concierge

ADD kernel-concierge $GOPATH/src/kernel-concierge
ADD Debuggin $GOPATH/src/kernel-concierge/Debuggin 
ADD Handlers $GOPATH/src/kernel-concierge/Handlers
ADD Jaeger $GOPATH/src/kernel-concierge/Jaeger
ADD Kafka $GOPATH/src/kernel-concierge/Kafka
ADD Pending $GOPATH/src/kernel-concierge/Pending
ADD Services $GOPATH/src/kernel-concierge/Services
ADD checkIsUp.sh $GOPATH/src/kernel-concierge

WORKDIR $GOPATH/src/kernel-concierge
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine

COPY --from=builder /go/src/kernel-concierge/main /app/
COPY --from=builder /go/src/kernel-concierge/checkIsUp.sh /app/

WORKDIR /app
CMD ["sh","checkIsUp.sh"]⏎