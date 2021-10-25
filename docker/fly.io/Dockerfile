FROM alpine

EXPOSE 8080
EXPOSE 8000
EXPOSE 60000/udp

RUN mkdir -p /opt/uhppoted
RUN mkdir -p /opt/uhppoted/html
RUN mkdir -p /etc/uhppoted
RUN mkdir -p /etc/uhppoted/httpd
RUN mkdir -p /etc/uhppoted/simulator
RUN mkdir -p /var/uhppoted/httpd/system

COPY uhppoted.conf /etc/uhppoted

COPY httpd/auth.json  /etc/uhppoted/httpd
COPY httpd/acl.grl    /etc/uhppoted/httpd
COPY httpd/system.grl /etc/uhppoted/httpd
COPY httpd/doors.grl  /etc/uhppoted/httpd
COPY httpd/cards.grl  /etc/uhppoted/httpd
COPY httpd/groups.grl /etc/uhppoted/httpd

COPY httpd/controllers.json /var/uhppoted/httpd/system
COPY httpd/doors.json       /var/uhppoted/httpd/system
COPY httpd/cards.json       /var/uhppoted/httpd/system
COPY httpd/groups.json      /var/uhppoted/httpd/system
COPY httpd/events.json      /var/uhppoted/httpd/system
COPY httpd/logs.json        /var/uhppoted/httpd/system

COPY simulator/201020304.json /etc/uhppoted/simulator
COPY simulator/303986753.json /etc/uhppoted/simulator
COPY simulator/405419896.json /etc/uhppoted/simulator

COPY uhppoted            /opt/uhppoted/
COPY uhppote-simulator   /opt/uhppoted/
COPY uhppoted-httpd      /opt/uhppoted/
COPY httpd/html/         /opt/uhppoted/html/
COPY httpd/translations/ /opt/uhppoted/translations/

WORKDIR /opt/uhppoted
RUN ["chmod", "+x", "uhppoted"]

ENTRYPOINT /opt/uhppoted/uhppoted