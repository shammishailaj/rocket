FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh-client ca-certificates

RUN adduser -D -g '' astrocorp

USER astrocorp

WORKDIR /rocket

ADD dist/rocket /bin/rocket

CMD ["rocket"]
