## MacOS Install Notes

### uhppoted-rest
```
sudo uhppoted-rest daemonize
```

Installs *uhppoted-rest* as a launchd system daemon:
- Creates launchd configuration file: `/Library/LaunchDaemons/com.github.twystd.uhppoted-rest.plist`
- Creates daemon working directory: `/usr/local/var/com.github.twystd.uhppoted`
- Adds *uhppoted-rest* executable to the application firewall (if enabled) and unblocks incoming connections

The default configuration is set to *run on load* and logs to the following files:
- `/usr/local/var/log/com.github.twystd.uhppoted-rest.log`
- `/usr/local/var/log/com.github.twystd.uhppoted-rest.err`

Start daemon:
```
sudo launchctl load /Library/LaunchDaemons/com.github.twystd.uhppoted-rest.plist
```

Stop daemon:
```
sudo launchctl unload /Library/LaunchDaemons/com.github.twystd.uhppoted-rest.plist
```