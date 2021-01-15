FROM alpine

EXPOSE 8080

RUN mkdir -p /etc/uhppoted
RUN mkdir -p /var/uhppoted

ADD uhppoted.conf /etc/uhppoted/
ADD ca.cert       /etc/uhppoted/rest/
ADD uhppoted.cert /etc/uhppoted/rest/
ADD uhppoted.key  /etc/uhppoted/rest/

WORKDIR /opt/uhppoted 
ADD uhppoted-rest     .
ADD uhppote-simulator .

ENTRYPOINT /opt/uhppoted/uhppoted-rest
