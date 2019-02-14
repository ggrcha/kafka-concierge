FROM golang:1.11.5-stretch as builder

RUN apt install git
#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
ADD . kernel/kernel-concierge
WORKDIR $GOPATH/kernel/kernel-concierge
RUN go mod download

# ADD kernel-concierge $GOPATH/src/kernel-concierge
# ADD Debuggin $GOPATH/src/kernel-concierge/Debuggin 
# ADD Handlers $GOPATH/src/kernel-concierge/Handlers
# ADD Jaeger $GOPATH/src/kernel-concierge/Jaeger
# ADD Kafka $GOPATH/src/kernel-concierge/Kafka
# ADD Pending $GOPATH/src/kernel-concierge/Pending
# ADD Services $GOPATH/src/kernel-concierge/Services
# ADD checkIsUp.sh $GOPATH/src/kernel-concierge

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine

COPY --from=builder /go/kernel/kernel-concierge/main /app/
COPY --from=builder /go/kernel/kernel-concierge/checkIsUp.sh /app/

WORKDIR /app
CMD ["sh","checkIsUp.sh"]⏎