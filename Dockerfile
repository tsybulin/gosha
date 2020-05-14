FROM alpine as build
RUN apk --no-cache add tzdata
RUN apk add -U --no-cache ca-certificates

FROM scratch as final
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Europe/Kiev
CMD ["/gosha"]
ADD gosha /
