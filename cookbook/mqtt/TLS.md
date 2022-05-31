# Creating the server and client keys and certificates for the test docker HiveMQ

- Ref. [HiveMQ: HowTo configure server-side TLS with HiveMQ and Keytool (self-signed)](https://www.hivemq.com/docs/hivemq/4.4/user-guide/howtos.html)

- Ref. [stackoverflow: How to add subject alernative name to ssl certs?](https://stackoverflow.com/questions/8744607/how-to-add-subject-alernative-name-to-ssl-certs#8744717)

### Create server keystore and certificates

1. Generate server key pair:

   keytool -genkey -keyalg RSA -alias hivemq -keystore localhost.jks -storepass hivemq -validity 365 -keysize 2048 -ext "SAN=DNS:localhost,IP:192.168.1.100"

>  What is your first and last name?
>    [Unknown]:  localhost
>  What is the name of your organizational unit?
>    [Unknown]:  hivemq
>  What is the name of your organization?
>    [Unknown]:  uhppoted
>  What is the name of your City or Locality?
>    [Unknown]:  docker
>  What is the name of your State or Province?
>    [Unknown]:  docker
>  What is the two-letter country code for this unit?
>    [Unknown]:  MQ

2. Export server certificate as a PEM file:

   keytool -exportcert -keystore localhost.jks -alias hivemq -keypass hivemq -storepass hivemq -rfc -file localhost.pem

3. Check the server certificate:

   openssl x509 -in localhost.pem -noout -text

4. Copy the server JKS file and certificates to the `docker` folders:

   cp localhost.jks ./docker/hivemq/localhost.jks
   cp localhost.pem ./docker/hivemq/localhost.pem
   cp localhost.pem ./docker/uhppoted-mqtt/hivemq.pem
   cp localhost.jks ./docker/integration-tests/hivemq/localhost.jks
   cp localhost.pem ./docker/integration-tests/hivemq/localhost.pem
   cp localhost.pem ./docker/integration-tests/mqttd/localhost.pem

5. Copy the server certificate to the `uhppoted` configuration folder:

   cp localhost.pem /usr/local/etc/com.github.uhppoted/hivemq.pem

### Create client keys and certificates

1. Generate client key pair without key password:

   openssl req -x509 -newkey rsa:2048 -keyout client.key -out client.cert -days 365 -nodes

>  Country Name (2 letter code) [AU]:MQ
>  State or Province Name (full name) [Some-State]:localhost
>  Locality Name (eg, city) []:localhost
>  Organization Name (eg, company) [Internet Widgits Pty Ltd]:localhost
>  Organizational Unit Name (eg, section) []:localhost
>  Common Name (e.g. server FQDN or YOUR name) []:localhost

2. Check the client certificate:

   openssl x509 -in client.cert -noout -text

3. Export client certificate as a DER file:

   openssl x509 -outform der -in client.cert -out client.crt

4. Add the client certificate to the server truststore:

   keytool -import -file client.crt -alias client -keystore clients.jks -storepass hivemq

5. Copy the client key and certificate to the `uhppoted` configuration folder:

   cp client.key  /usr/local/etc/com.github.uhppoted/mqtt/client.key
   cp client.cert /usr/local/etc/com.github.uhppoted/mqtt/client.cert

6. Copy the truststore, client key and client certificate to the `docker` folders:

   cp clients.jks ~/Development/uhppote/uhppoted/docker/hivemq/clients.jks
   cp clients.jks ~/Development/uhppote/uhppoted/docker/integration-tests/hivemq/clients.jks
   cp client.key  ~/Development/uhppote/uhppoted/docker/hivemq/client.key
   cp client.key  ~/Development/uhppote/uhppoted/docker/uhppoted-mqtt/secure/client.key
   cp client.key  ~/Development/uhppote/uhppoted/docker/integration-tests/hivemq/client.key
   cp client.key  ~/Development/uhppote/uhppoted/docker/integration-tests/mqttd/secure/client.key
   cp client.cert ~/Development/uhppote/uhppoted/docker/hivemq/client.cert
   cp client.cert ~/Development/uhppote/uhppoted/docker/uhppoted-mqtt/secure/client.cert
   cp client.cert ~/Development/uhppote/uhppoted/docker/integration-tests/hivemq/client.cert
   cp client.cert ~/Development/uhppote/uhppoted/docker/integration-tests/mqttd/secure/client.cert
   cp client.crt  ~/Development/uhppote/uhppoted/docker/hivemq/client.crt
   cp client.crt  ~/Development/uhppote/uhppoted/docker/integration-tests/hivemq/client.crt


### Rebuild the Docker images

   cd ./docker/hivemq        && docker build -f Dockerfile -t hivemq/uhppoted . 
   cd ./docker/uhppoted-mqtt && docker build -f Dockerfile -t uhppoted/mqtt   . 

