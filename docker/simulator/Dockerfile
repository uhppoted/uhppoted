FROM alpine

EXPOSE 8000
EXPOSE 60000/udp

RUN mkdir -p /opt/uhppote
RUN mkdir -p /etc/uhppote/devices
ADD ./docker/simulator/devices/303986753.json /etc/uhppote/devices
ADD ./docker/simulator/devices/405419896.json /etc/uhppote/devices

WORKDIR /opt/uhppote
ADD ./docker/simulator/uhppote-simulator .

ENTRYPOINT /opt/uhppote/uhppote-simulator --debug --devices /etc/uhppote/devices
