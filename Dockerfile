FROM alpine
LABEL name=ddns-ipv6-go
LABEL url=https://github.com/penggege18/ddns-ipv6-go
RUN apk add --no-cache curl grep

WORKDIR /app
COPY ddns-ipv6-go /app/
COPY zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
EXPOSE 9876
ENTRYPOINT ["/app/ddns-ipv6-go"]
CMD ["-l", ":9876", "-f", "300"] 