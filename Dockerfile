FROM golang:1.22.6 AS builder
WORKDIR /
COPY . .
RUN go build -o bot

FROM apache/tika:2.9.2.1-full AS runner
COPY --from=builder ./bot ./

ENV TIKA_HOST=http://localhost:9998

# calling bot then default entrypoint of tika 2.9.2.1-full
ENTRYPOINT ["/bin/sh", "-c", "./bot & exec java -cp \"/tika-server-standard-${TIKA_VERSION}.jar:/tika-extras/*\" org.apache.tika.server.core.TikaServerCli -h 0.0.0.0 $0 $@"]
