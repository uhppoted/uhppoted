FROM alpine

EXPOSE 8080

RUN mkdir -p /etc/uhppoted
ADD ./docker/uhppoted-rest/uhppoted.conf /etc/uhppoted/
ADD ./docker/uhppoted-rest/ca.cert       /etc/uhppoted/rest/
ADD ./docker/uhppoted-rest/uhppoted.cert /etc/uhppoted/rest/
ADD ./docker/uhppoted-rest/uhppoted.key  /etc/uhppoted/rest/

RUN mkdir -p /var/uhppoted

WORKDIR /opt/uhppoted 
ADD ./docker/uhppoted-rest/uhppoted-rest     .
ADD ./docker/uhppoted-rest/uhppote-simulator .

ENTRYPOINT /opt/uhppoted/uhppoted-rest
