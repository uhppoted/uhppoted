## MacOS Install Notes

1. If Application Firewall is enabled, incoming UDP is blocked. Pending implementation of socket handoff, 
a partial workaround is to add uhppoted to the firewall:
```
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate off
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/bin/uhppoted
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /usr/local/bin/uhppoted
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on
```
This seems to be required on every startup.
