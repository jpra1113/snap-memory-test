FROM alpine:3.4

EXPOSE 8181

COPY snap-memory-test snap-memory-test

ENTRYPOINT ["./snap-memory-test"]
