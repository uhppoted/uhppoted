ARG TAG=latest

FROM hivemq/hivemq4:${TAG} 

COPY --chown=hivemq:hivemq ./docker/hivemq/config.xml /opt/hivemq/conf/config.xml
COPY --chown=hivemq:hivemq ./docker/hivemq/localhost.jks /opt/hivemq/conf/hivemq.jks

CMD /opt/hivemq/bin/run.sh
