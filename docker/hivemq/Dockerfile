ARG TAG=latest

FROM hivemq/hivemq4:${TAG} 

COPY --chown=hivemq:hivemq config.xml    /opt/hivemq/conf/config.xml
COPY --chown=hivemq:hivemq localhost.jks /opt/hivemq/conf/hivemq.jks
COPY --chown=hivemq:hivemq clients.jks   /opt/hivemq/conf/clients.jks

CMD /opt/hivemq/bin/run.sh
