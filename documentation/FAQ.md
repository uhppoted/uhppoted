# FAQ

1. _[Upgrading UT0311-L0x Firmware](https://github.com/uhppoted/uhppoted/edit/main/documentation/FAQ.md#1-upgrading-ut0311-l0x-firmware)_
2. _[Ephemeral ports and binding to `0.0.0.0:0`](https://github.com/uhppoted/uhppoted/edit/main/documentation/FAQ.md#2-ephemeral-ports-and-binding-to-00000)_
3. _[Docker + UDP](https://github.com/uhppoted/uhppoted/edit/main/documentation/FAQ.md#3-docker--udp)_
4. _[Entering card numbers on a keypad](https://github.com/uhppoted/uhppoted/blob/main/documentation/FAQ.md#4-entering-card-numbers-on-a-keypad)_
5. _[`bind`, `broadcast`, `listen` and `listener` addresses](https://github.com/uhppoted/uhppoted/blob/main/documentation/FAQ.md#5-bind-broadcast-listen-and-listener-addresses)_


----
### 1. Upgrading UT0311-L0x Firmware

**NOTE: the information below is out of date - apparently more recent versions can be upgraded using the manufacturer application.**

_(cf. [Upgrading UT0311-L0x Firmware](https://github.com/uhppoted/uhppote-core/issues/6))_

According to the vendor's [Amazon page](https://www.amazon.com/UHPPOTE-Professional-Wiegand-Network-Software):

**Q: How do I upgrade the firmware on the boards?**\
A: I am sorry the firmware can't been upgraded. thanks.\
_By Xiaojia Huang on September 9, 2019_

**Q: Is there a way to update the firmware or software**\
A: the software have not updated sevral year, so don't need see more\
_By Xiaojia Huang on April 14, 2017_

**Q: Are you able to upgrade the firmware on this ? if check is made on the door using the software? how can you do this as i want to buy more**\
A: Here are the anwsers of your questions:\
1\)The firmware can't be upgraded, and you could download the latest software.\
2\) what kind of the check do you want to made on the door using the software. We recommend that you could tell us. Thanks.\
_By Xiaojia Huang on January 18, 2022_

----

### 2. Ephemeral ports and binding to `0.0.0.0:0`

As per [Microsoft Knowledgebase Article 929851](https://learn.microsoft.com/en-us/troubleshoot/windows-server/networking/default-dynamic-port-range-tcpip-chang),
the default Windows ephemeral port range extends from 49152 to 65535, which includes the default UHPPOTE UDP port (`60000`). Present-day MacOS, BSD and Linux
have similar ranges.

If an application is assigned port `60000` when binding to `0.0.0.0:0` it will receive any outgoing UDP broadcast requests and interpret
them as replies - which will be, uh, more than a little confusing, e.g.:
```
request:
   17 94 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00

reply:
   17 94 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00

reply:
   17 94 00 00 78 37 2a 18  c0 a8 01 64 ff ff ff 00
   c0 a8 01 01 00 12 23 34  45 56 08 92 20 18 11 05
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00
      
get-all-controllers:
   controller: 0
      address: 0.0.0.0
      netmask: 0.0.0.0
      gateway: 0.0.0.0
          MAC: 00:00:00:00:00:00
      version: v0.00
         date: ---

   controller: 405419896
      address: 192.168.1.100
      netmask: 255.255.255.0
      gateway: 192.168.1.1
          MAC: 00:12:23:34:45:56
      version: v8.92
         date: 2018-11-05
```

Barring a programming or configuration error which assigns a bind port of 60000, this condition has an about 1 in ~20000 chance of occurring and is unlikely
to be picked up during testing. The _uhppoted_ libraries now return an error (or raise exception) when the bind port is the same as the broadcast port, and
the application is expected to retry.

However, the correct fix is to exclude port 60000 from the ephemeral ports that can be assigned by the OS:

#### Linux

Ref. [Understanding and Utilizing Linux Ephemeral Port Range](https://linuxvox.com/blog/linux-ephemeral-port-range).

1. Show the current ephemeral port range and reserved ports:
```
sysctl net.ipv4.ip_local_port_range 
sysctl net.ipv4.ip_local_reserved_ports 
```

2. Limit the ephemeral port range to exclude port 60000, e.g.:
```
sudo sysctl -w net.ipv4.ip_local_port_range="49152 59999"
```
or
```
sudo sysctl -w net.ipv4.ip_local_reserved_ports=60000"
```

3. Update _/etc/sysctl.conf_ with:
```
net.ipv4.ip_local_port_range = 49152 59999
```
or
```
net.ipv4.ip_local_reserved_ports=60000
```

4. _(optionally)_ reload configuration:
```
sudo sysctl -p
```

#### MacOS

**NOT TESTED - USE AT OWN RISK**

See [com.startup.sysctl.plist](https://gist.github.com/pythoninthegrass/8073e5e3b24f385c9d9b712f6f243982). 

1. Show the current ephemeral port range:
```
sysctl net.inet.ip.portrange.first net.inet.ip.portrange.last
```
2. Limit the ephemeral port range to exclude port 60000, e.g.:
```
sudo sysctl -w net.inet.ip.portrange.first=49152
sudo sysctl -w net.inet.ip.portrange.last=59999
```
3. Create a _launchd_ daemon:
    - _/Library/LaunchDaemons/com.user.sysctl.plist_
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.sysctl</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/sbin/sysctl</string>
        <string>-w</string>
        <string>net.inet.ip.portrange.first=49152</string>
        <string>net.inet.ip.portrange.last=59999</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
```
4. Reload configuration:
```
sudo launchctl load -w /Library/LaunchDaemons/com.user.sysctl.plist
```
or
```
sudo launchctl bootstrap system /Library/LaunchDaemons/com.user.sysctl.plist
```

#### Windows

1. Show the current ephemeral port range:
```
netsh int ipv4 show dynamicport udp
```
2. Limit the ephemeral port range to exclude port 60000, e.g.:
```
netsh int ipv4 set dynamicport udp start=49152 num=10848
```
or
```
netsh int ipv4 add excludedportrange protocol=udp startport=60000 numberofports=1
```

#### Other
For systems where it's not possible to modify the ephemeral port range, other workarounds include:
- If you know the destination address use the TCP transport protocol.
- Bind a _netcat_ or _socat_ listener to port `60000` before running the application:
   ```
   nc -lu 60000 > /dev/null
   ```
   or
   ```
   socat -u UDP-RECV:60000 - > /dev/null
   ```

#### References:
1. [_The Ephemeral Port Range_](https://www.ncftp.com/ncftpd/doc/misc/ephemeral_ports.html)
2. [_How to change/view the ephemeral port range on Windows machines?_](https://stackoverflow.com/questions/7006939/how-to-change-view-the-ephemeral-port-range-on-windows-machines#7007159)
3. [_You cannot exclude ports by using the ReservedPorts registry key in Windows Server 2008 or in Windows Server 2008 R2_](https://support.microsoft.com/en-us/topic/you-cannot-exclude-ports-by-using-the-reservedports-registry-key-in-windows-server-2008-or-in-windows-server-2008-r2-a68373fd-9f64-4bde-9d68-c5eded74ea35)
4. [_Listen to UDP data on local port with netcat_](https://serverfault.com/questions/207683/listen-to-udp-data-on-local-port-with-netcat)
5. [How to view and edit the ephemeral port range on Linux?](https://stackoverflow.com/questions/28573390/how-to-view-and-edit-the-ephemeral-port-range-on-linux)
6. [Understanding and Utilizing Linux Ephemeral Port Range](https://linuxvox.com/blog/linux-ephemeral-port-range)
7. [How safe is it to change the Linux Ephemeral Port range](https://unix.stackexchange.com/questions/211647/how-safe-is-it-to-change-the-linux-ephemeral-port-range)
8. [SYSCTL Command: Get or set the kernel state in macOS](https://ss64.com/mac/sysctl.html)
9. [com.startup.sysctl.plist](https://gist.github.com/pythoninthegrass/8073e5e3b24f385c9d9b712f6f243982)
10. [Powershell: Check Excluded Port Ranges](https://powershell.xevion.dev/scripts/excluded-ports)


----
### 3. Docker + UDP

In _bridge networking mode_ (_MacOS_ and _Windows_), the Docker UDP transport drops incoming packets at a significantly higher rate than 
typically experienced on a LAN/WAN. _host networking mode_ (Linux, RaspberryPi, etc) appears to operate normally and reliably.

The TCP transport doesn't appear to be affected.

----
### 4. Entering card numbers on a keypad

_(cf. [Entering card numbers on a keypad](https://github.com/uhppoted/uhppoted-lib-python/discussions/14#discussioncomment-15113138))_

For systems with a keypad but without a card reader it is possible to use card numbers for access by entering:
```
*card-number* or *card-number*PIN*
```
e.g.:
```
*10058400* or *10058400*7531*
```

If you're using Windows app, configure the controller with:
```
Configuration > Password Management > Card + PIN
Configuration > Password Management > Manual Input Password
```

----
### 5. `bind`, `broadcast`, `listen` and `listener` addresses

#### `bind` address

The `bind` address is the IPv4 address of the host computer (RaspberryPi, etc, etc). For most configurations, using "0.0.0.0"
(_INADDR_ANY_) is fine, but you may want to bind it to a specific interface if you have:
   - a machine with multiple network cards
   - are running both WiFi and a wired LAN
   - a firewall rule for outgoing connections that restrict the source port


#### `broadcast` address

The `broadcast` address is the IPv4 UDP broadcast address of the host computer (RaspberryPi, etc, etc). For most configurations, 
using "255.255.255.255" is fine, but you may want to use to use the `broadcast` address of a specific interface if you have:
   - a machine with multiple network cards
   - are running both WiFi and a wired LAN


#### `listen` address

The `listen` address is the IPv4 address of the host computer (RaspberryPi, etc, etc) on which to receive events from a controller.
For most configurations, something like "0.0.0.0:60001" is fine, but you may want to bind it to a specific interface if you have:
   - a machine with multiple network cards
   - are running both WiFi and a wired LAN
   - a firewall rule for incoming UDP messages connections that restrict the destination interface

The _unhppoted_xxx_ subsystems default to UDP port 60001 but you may use any port you like.


#### `listener` address

The `listener` address is the IPv4 address to which the controller sends UDP event messages i.e. it is the IPv4 _address:port_ of the 
host computer (RaspberryPi, etc, etc) which expects to receive events from a controller.
  
#### Example

Assuming:
- a host computer with IPv4 address 192.169.1.100
- a UHPPOTE controller with IPv4 address 192.168.1.25
- you are using UDP port 60001 for events

then:

- the `bind` address is either `0.0.0.0` or `192.168.1.100`
- the `broadcast` address is (typically) `255.255.255.255` or `192.168.1.255`
- the `listen` address is either `0.0.0.0:60001` or `192.168.1.100:60001`
- the `listener` address for the _set-listener_ function is `192.168.1.100:60001`
