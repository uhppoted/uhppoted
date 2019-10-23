## MacOS Install Notes

### uhppoted
```
sudo uhppoted daemonize
```

Installs *uhppoted* as a launchd system daemon:
- Creates launchd configuration file: `/Library/LaunchDaemons/com.github.twystd.uhppoted.plist`
- Creates daemon working directory: `/usr/local/var/com.github.twystd.uhppoted`
- Adds *uhppoted* executable to the application firewall (if enabled) and unblocks incoming connections

The default configuration is set to *run on load* and logs to the following files:
- `/usr/local/var/log/com.github.twystd.uhppoted.log`
- `/usr/local/var/log/com.github.twystd.uhppoted.err`

Start daemon:
```
sudo launchctl load /Library/LaunchDaemons/com.github.twystd.uhppoted.plist
```

Stop daemon:
```
sudo launchctl unload /Library/LaunchDaemons/com.github.twystd.uhppoted.plist
```