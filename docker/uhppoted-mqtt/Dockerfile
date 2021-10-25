FROM alpine

RUN mkdir -p /etc/uhppoted
RUN mkdir -p /etc/uhppoted/mqtt
RUN mkdir -p /etc/uhppoted/mqtt/rsa
RUN mkdir -p /etc/uhppoted/mqtt/rsa/signing
RUN mkdir -p /etc/uhppoted/mqtt/rsa/encryption

RUN mkdir -p /var/uhppoted

COPY uhppoted.conf                   /etc/uhppoted
COPY hivemq.pem                      /etc/uhppoted
COPY secure/client.key               /etc/uhppoted/mqtt
COPY secure/client.cert              /etc/uhppoted/mqtt
COPY secure/rsa/signing/mqttd.key    /etc/uhppoted/mqtt/rsa/signing
COPY secure/rsa/encryption/mqttd.key /etc/uhppoted/mqtt/rsa/encryption

WORKDIR /opt/uhppoted 
COPY     uhppoted-mqtt .

ENTRYPOINT /opt/uhppoted/uhppoted-mqtt --console
