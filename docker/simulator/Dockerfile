FROM alpine

EXPOSE 8000
EXPOSE 60000/udp

RUN mkdir -p /opt/uhppote
RUN mkdir -p /etc/uhppote/devices
ADD ./docker/simulator/devices/33333.json     /etc/uhppote/devices
ADD ./docker/simulator/devices/305419896.json /etc/uhppote/devices

WORKDIR /opt/uhppote
ADD ./docker/simulator/uhppote-simulator .

ENTRYPOINT /opt/uhppote/uhppote-simulator --debug --devices /etc/uhppote/devices
