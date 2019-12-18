ARG TAG=latest

FROM hivemq/hivemq4:${TAG} 

COPY --chown=hivemq:hivemq ./docker/hivemq/config.xml /opt/hivemq/conf/config.xml
COPY --chown=hivemq:hivemq ./docker/hivemq/localhost.jks /opt/hivemq/conf/hivemq.jks

# COPY --chown=hivemq:hivemq myentrypoint.sh /opt/myentrypoint.sh
#
# RUN chmod +x /opt/myentrypoint.sh 
#
# ENTRYPOINT ["/opt/myentrypoint.sh"] 
