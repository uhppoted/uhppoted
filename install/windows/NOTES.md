## Windows Install Notes

### uhppoted
```
sudo uhppoted daemonize
```

Installs *uhppoted* as a Windows service:
- Registers uhppoted as a Windows service

Start service:
```
net start uhppoted
```

Stop service:
```
net start uhppoted
```

Query service state:
```
sc queryex uhppoted
```

Kill service
```
taskkill /f /pid <PID>
```

Test URL's:
- http://127.0.0.1:8001/uhppote/device
- http://127.0.0.1:8001/uhppote/device/405419896/time