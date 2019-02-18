FROM golang:1.11.5-stretch as builder

# RUN apt install git
#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first

RUN mkdir /kernel-concierge
WORKDIR /kernel-concierge

ADD go.mod .
ADD go.sum .

RUN go mod download

ADD . .

# ADD kernel-concierge $GOPATH/src/kernel-concierge
# ADD Debuggin $GOPATH/src/kernel-concierge/Debuggin 
# ADD Handlers $GOPATH/src/kernel-concierge/Handlers
# ADD Jaeger $GOPATH/src/kernel-concierge/Jaeger
# ADD Kafka $GOPATH/src/kernel-concierge/Kafka
# ADD Pending $GOPATH/src/kernel-concierge/Pending
# ADD Services $GOPATH/src/kernel-concierge/Services
# ADD checkIsUp.sh $GOPATH/src/kernel-concierge

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/kernel-concierge .

FROM alpine

COPY --from=builder /go/bin/kernel-concierge /app/
COPY checkIsUp.sh /app/

WORKDIR /app
CMD ["sh","checkIsUp.sh"]⏎