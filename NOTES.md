## WORKING NOTES

#### Linux

1. UDP broadcast on Ubuntu needs the following UFW rule:
   - ufw allow from <local address> to any port 60000 proto udp

#### MacOS

1. Out of the box, MacOS doesn't support UDP broadcast on the loopback interface. Binding to 
   INADDR_ANY binds to the actual interface and seems to work ok for use with *uhppoted* and
   the *simulation*.

#### Windows

1. [/dev/null](https://stackoverflow.com/questions/313111/is-there-a-dev-null-on-windows)

#### Docker

1. [UDP port forwarding](https://stackoverflow.com/questions/42422406/receive-udp-multicast-in-docker-container)

### udp-broadcast-relay

1. https://github.com/nomeata/udp-broadcast-relay
2. https://github.com/udp-redux/udp-broadcast-relay-redux
3. https://forum.opnsense.org/index.php?topic=11818.0
4. https://networkengineering.stackexchange.com/questions/71202/how-to-route-incoming-udp-unicast-traffic-to-multiple-computers

### uhppoted

#### Certificate Authority

- [CloudFlare CFSSL](https://blog.cloudflare.com/introducing-cfssl/)
- [github:CloudFlare CFSSL](https://github.com/cloudflare/cfssl)
- [OpenSSL Certificate Authority](https://jamielinux.com/docs/openssl-certificate-authority/index.html)
- Chrome requires a Subject Alternative Address of the form DNS:<hostname>,IP:<IP address>
- [ServerFault:ERR_CERT_COMMON_NAME_INVALID](https://serverfault.com/questions/880804/can-not-get-rid-of-neterr-cert-common-name-invalid-error-in-chrome-with-self)
- If using the [OpenSSL Certificate Authority](https://jamielinux.com/docs/openssl-certificate-authority/index.html)guide, 
  update the *server_cert* section of the intermediate CA openssl.cnf with
```
  [ server_cert ]
  ...
  ...
  subjectAltName = DNS:<hostname>,IP:<IP address>
```

#### Encryption/Signing

The current MQTTD sign-then-encrypt implementation knowingly allows for *surreptitious forwarding* - MQTT 3.1x does not provide
a way to identify the actual sender of a message. This does somewhat impact system security e.g.:

- a geo-fenced access control system that requires a user to be present in an area to open a door can be co-operatively 
  circumvented if the authorised user not inside the geo-fenced area sends a signed 'OPEN' request to a 
  non-authorised user inside the geo-fenced area who then forwards it to the access control system.

## CRDTs

- [Bidirectional Program Transformations](https://youtu.be/1gGd7pKSpRM?t=1855)

References:

1. https://crypto.stackexchange.com/questions/8139/secure-encrypt-then-sign-with-rsa
2. http://world.std.com/~dtd/sign_encrypt/sign_encrypt7.html
3. https://askubuntu.com/questions/1093591/how-should-i-change-encryption-according-to-warning-deprecated-key-derivat
4. https://superuser.com/questions/1016696/using-a-hash-other-than-sha1-for-oaep-with-openssl-cli
5. https://security.stackexchange.com/questions/185083/specifying-rsa-oaep-label-via-openssl-command-line
6. https://crypto.stackexchange.com/questions/202/should-we-mac-then-encrypt-or-encrypt-then-mac

