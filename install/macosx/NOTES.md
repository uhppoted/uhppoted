## MacOS Install Notes

### uhppoted
```
sudo uhpppoted install
```

Installs *uhppoted* as a launchd system daemon:
- Creates launchd configuration file */Library/LaunchDaemons/com.github.twystd.uhppoted.plist*
- Creates the working directory */usr/local/var/com.github.twystd.uhppoted*
- Adds *uhppoted* to the application firewall (if enabled) and unblocks incoming connections

The default configuration is set to 'run on load' and logs to the following files:
- `/usr/local/var/log/com.github.twystd.uhppoted.log`
- `/usr/local/var/log/com.github.twystd.uhppoted.err`

Start the daemon:
```
sudo launchctl load /Library/LaunchDaemons/com.github.twystd.uhppoted.plist
```

Stop the daemon:
```
sudo launchctl unload /Library/LaunchDaemons/com.github.twystd.uhppoted.plist
```