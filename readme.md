My laptop starts overheating if I run the official [dropbox daemon](https://www.dropbox.com/install?os=lnx).

Run this hand-rolled dropbox sync script in a cron job as a more lightweight solution.

```
$ crontab -e
```

Add:
```
@daily /home/andrew/development/go/bin/dropbox-nightly-sync
```
