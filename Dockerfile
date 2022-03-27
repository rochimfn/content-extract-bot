FROM adoptopenjdk/openjdk8:alpine-slim
COPY . ./
RUN chmod +x ./gradlew
RUN ./gradlew build

FROM adoptopenjdk/openjdk8:alpine-slim
WORKDIR /build
RUN apk update && apk --no-cache add tesseract-ocr
COPY --from=0 /app/build/distributions/app.tar ./
RUN tar -xf app.tar && rm app.tar
CMD ["/build/app/bin/app"]
