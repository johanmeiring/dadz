FROM alpine

RUN apk add --update ca-certificates

COPY build/dadz /dadz

ENTRYPOINT ["/dadz"]
