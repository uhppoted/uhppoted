ARG TAG=latest

FROM hivemq/hivemq4:${TAG} 

# ENV MY_CUSTOM_PLUGIN_ENV myvalue 
# ENV HIVEMQ_CLUSTER_PORT 8000

COPY --chown=hivemq:hivemq ./docker/hivemq/config.xml /opt/hivemq/conf/config.xml
COPY --chown=hivemq:hivemq ./docker/hivemq/hivemq.jks /opt/hivemq/conf/hivemq.jks

# COPY --chown=hivemq:hivemq your-license.lic /opt/hivemq/license/your-license.lic 
# COPY --chown=hivemq:hivemq myplugin.jar /opt/hivemq/plugins/myplugin.jar
# COPY --chown=hivemq:hivemq myentrypoint.sh /opt/myentrypoint.sh
#
# RUN chmod +x /opt/myentrypoint.sh 
#
# ENTRYPOINT ["/opt/myentrypoint.sh"] 