## WORKING NOTES

### UDP broadcast

#### Linux

1. UDP broadcast on Linux between *uhppoted* and *simulation* only works for the loopback interface. 
   Binding to INADDR_ANY seems to be sufficient.

2. The *simulation* does not receive any UDP packets when *uhppoted* is bound to the interface IP and 
   is using the interface broadcast address.

3. Other things to try:
   - https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc
   - https://developerweb.net/viewtopic.php?id=5722
   - https://github.com/golang/go/issues/6935

4. Suggestions that didn't work:
   - DialUDP rather than ListenUDP
   - Set SO_BROADCAST on listening socket

#### MacOS

1. Out of the box, MacOS doesn't support UDP broadcast on the loopback interface. Binding to 
   INADDR_ANY binds to the actual interface and seems to work ok for use with *uhppoted* and
   the *simulation*.

