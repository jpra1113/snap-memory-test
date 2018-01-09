FROM ubuntu:16.04

EXPOSE 8181

COPY snap-memory-test snap-memory-test

ENTRYPOINT ["./snap-memory-test"]
