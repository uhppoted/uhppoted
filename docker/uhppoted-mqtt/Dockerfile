FROM alpine

RUN mkdir -p /etc/uhppoted
RUN mkdir -p /etc/uhppoted/mqtt
RUN mkdir -p /etc/uhppoted/mqtt/rsa
RUN mkdir -p /etc/uhppoted/mqtt/rsa/signing
RUN mkdir -p /etc/uhppoted/mqtt/rsa/encryption

RUN mkdir -p /var/uhppoted

ADD uhppoted.conf                   /etc/uhppoted
ADD hivemq.pem                      /etc/uhppoted
ADD secure/client.key               /etc/uhppoted/mqtt
ADD secure/client.cert              /etc/uhppoted/mqtt
ADD secure/rsa/signing/mqttd.key    /etc/uhppoted/mqtt/rsa/signing
ADD secure/rsa/encryption/mqttd.key /etc/uhppoted/mqtt/rsa/encryption

WORKDIR /opt/uhppoted 
ADD     uhppoted-mqtt .

ENTRYPOINT /opt/uhppoted/uhppoted-mqtt --console
