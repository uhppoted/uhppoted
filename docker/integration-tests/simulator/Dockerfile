FROM alpine

EXPOSE 8000
EXPOSE 60000/udp

RUN mkdir -p /opt/uhppote
RUN mkdir -p /etc/uhppote/devices/60000
RUN mkdir -p /etc/uhppote/devices/60001
ADD devices/303986753.json /etc/uhppote/devices/60000
ADD devices/405419896.json /etc/uhppote/devices/60000
ADD devices/301357913.json /etc/uhppote/devices/60001
ADD devices/402468024.json /etc/uhppote/devices/60001

WORKDIR /opt/uhppote
ADD uhppote-simulator .

ENTRYPOINT /opt/uhppote/uhppote-simulator --debug --devices /etc/uhppote/devices/60000
