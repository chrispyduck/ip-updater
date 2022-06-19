FROM registry1.dso.mil/ironbank/google/golang/golang-1.18 AS base

FROM base AS build 
WORKDIR /work 
COPY * /work/
RUN --mount=type=cache,target=/go,uid=1001,gid=1001 go get -d -v 
RUN --mount=type=cache,target=/tmp,uid=1001,gid=1001 go build -v .

FROM base AS RUN
COPY --from=build /work/ip-updater /bin
ENTRYPOINT ["ip-updater"]