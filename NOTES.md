## WORKING NOTES

#### Linux

1. UDP broadcast on Ubuntu needs the following UFW rule:
   - ufw allow from <local address> to any port 60000 proto udp

#### MacOS

1. Out of the box, MacOS doesn't support UDP broadcast on the loopback interface. Binding to 
   INADDR_ANY binds to the actual interface and seems to work ok for use with *uhppoted* and
   the *simulation*.

#### Windows

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
