FROM gradle:7.4.1-jdk8
WORKDIR /build
COPY . ./
RUN ./gradlew build

FROM apache/tika:1.28.1-full
WORKDIR /app
COPY --from=0 /build/app/build/distributions/app.tar ./
RUN tar -xf app.tar && rm app.tar
CMD ["/app/app/bin/app"]
