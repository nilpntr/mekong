FROM sammobach/go:1.16 as build
LABEL maintainer="Sam Mobach <hello@sammobach.com>"
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN mkdir -p /app/dist && go build -ldflags "-s -w" -o dist/mekong ./cmd/mekong
CMD ["/app/dist/mekong"]

FROM alpine:3.14.0
LABEL maintainer="Sam Mobach <hello@sammobach.com>"
RUN mkdir /app
COPY --from=build /app/dist /app
WORKDIR /app
CMD ["/app/mekong"]
