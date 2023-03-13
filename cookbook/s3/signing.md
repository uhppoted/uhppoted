# Creating a signed ACL file

#### Create RSA signing keys

```
   openssl genrsa -out QWERTY54.key 2048
   openssl rsa    -in  QWERTY54.key -out QWERTY54.pub -outform PEM -pubout
```

#### Copy the public signing key to the `s3` configuration directory

```
   cp QWERTY54.pub /usr/local/etc/com.github.uhppoted/s3/rsa/signing/QWERTY54.pub
```

#### Sign the ACL file with the private key

```
   openssl dgst -sha256 -sign QWERTY54.key -out signature hogwarts.acl
```

#### Package the ACL and signature as a .tar.gz file

Create a _.tar.gz_ file containg the _ACL_ and _signature_ files, with the `uname` and `gname` set to the
user name of the key used to sign the ACL file:

```
   tar --uname=QWERTY54 --gname=QWERTY54 -cvzf hogwarts.tar.gz hogwarts.acl signature
```
