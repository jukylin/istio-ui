FROM alpine:latest
RUN apk add bash
ADD ./istio-ui /
ADD ./views/ /views/
ADD ./conf/* /conf/
RUN mkdir -p /data/www/istio_config
RUN mkdir -p /data/www/istio_upload
EXPOSE 9100
CMD ["/istio-ui"]
