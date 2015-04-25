# Deployment

Deploy with `bin/deploy`.

This lives on a regular server, not heroku.

Binary and comments get untarred into `/usr/local/rikki`:

```
/usr/local/rikki/rikki
/usr/local/rikki/comments/*
```

Logs can be found in `/var/log/upstart/rikki.log`.

## Environment variables.

RIKKI_SECRET has to match the one in the exercism.io application that is running
EXERCISM has to match the url of the exercism.io application
ANALYSSEUR has to match the url of the analysseur application that is running
REDIS has is where the jobs on queue 'analyze' are being taken from

## Upstart script

Config lives in `/etc/init/rikki.conf`

```
description "rikki"
author "Katrina Owen"

start on filesystem
stop on runlevel [!2345]

respawn

script
export REDIS=<redis url>
export EXERCISM=http://exercism.io
export ANALYSSEUR=http://analysseur.exercism.io
export RIKKI_SECRET=<shared secret>
export RIKKI_FEEDBACK_DIR=/usr/local/rikki/current/comments
/usr/local/rikki/current/rikki -exercism=$EXERCISM -analysseur=$ANALYSSEUR -redis=$REDIS
end script
```

Stop and start with:

```
service rikki start
service rikki stop
```
