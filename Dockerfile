FROM golang:alpine AS builder

ADD . /app
WORKDIR /app
RUN go build -o platform-info .
RUN chmod +x platform-info

#create new clean image
FROM alpine
WORKDIR /app
ADD clis/* /usr/bin/
RUN chmod +x /usr/bin/pks
RUN chmod +x /usr/bin/govc
RUN chmod +x /usr/bin/bosh

COPY --from=builder /app/platform-info /app/platform-info
RUN cp platform-info /usr/bin/platform-info