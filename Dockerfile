FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh-client ca-certificates

ADD dist/rocket /bin/rocket

RUN adduser -D -g '' astrocorp

USER astrocorp

WORKDIR /rocket

CMD ["rocket"]
