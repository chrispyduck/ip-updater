FROM golang:alpine AS build
WORKDIR /work 
COPY * /work/
RUN --mount=type=cache,target=/go go get -d -v 
RUN --mount=type=cache,target=/tmp go build -v -ldflags="-w -s" .
#RUN mkdir /out 
#RUN mv -v /work/ip-updater /out/

FROM alpine:latest AS run
WORKDIR /work
COPY --from=build /work/ip-updater /bin/ip-updater
ENTRYPOINT ["ip-updater"]