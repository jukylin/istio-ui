FROM alpine:latest
ADD ./istio-ui /
ADD ./views/ /views/
ADD ./conf/* /conf/
RUN mkdir -p /data/www/istio_config
RUN mkdir -p /data/www/istio_upload
RUN mkdir -p /data/www/istio_back_up
EXPOSE 9100
CMD ["/istio-ui"]
