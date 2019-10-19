## WORKING NOTES

### UDP broadcast

#### Linux

1. UDP broadcast on Ubuntu needs the following UFW rules:
   - ?
   - ?

#### MacOS

1. Out of the box, MacOS doesn't support UDP broadcast on the loopback interface. Binding to 
   INADDR_ANY binds to the actual interface and seems to work ok for use with *uhppoted* and
   the *simulation*.

