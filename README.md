Mayhem
===

A fork of [mayhem](https://github.com/BOTbkcd/mayhem) with some cut features
and generally simplified layout. It uses a local, updated, and modified version of the
[pure Go sqlite driver](https://github.com/glebarez/sqlite)

- No recurring tasks
- No steps
- Control over database storage
- Auto-saving to JSON
- Import/export via JSON

[![build](https://github.com/enckse/mayhem/actions/workflows/build.yml/badge.svg)](https://github.com/enckse/mayhem/actions/workflows/build.yml)

## usage

### configuration

`mayhem` is configured via TOML (in `MAYHEM_CONFIG`, `XDG_CONFIG_HOME`, or `$HOME/.config/mayhem/`) in a settings.toml file

```
[data]
# override the location where data is stored
directory="~/.mayhem"

[backups]
# enable backups into a directory (offset from data.directory)
# backups are taken when mayhem starts
directory="backups"
# keeps this duration of backups
duration="72h"
# control the format of the date, allows controlling how many backups one gets
format="20060102"

[json]
# save a JSON export on exit
exit=true
# save a JSON export on save/sync
sync=true
```

### usage

Run `mayhem` and follow the navigation keys/help

#### import/export

To import (or use `merge` instead to load into an existing mayhem database)
```
cat my_tasks.json | mayhem import
```

To export (save result as JSON syncing)
```
mayhem export > my_tasks.json
```

## build

clone and `make`

[![build](https://github.com/enckse/mayhem/actions/workflows/build.yml/badge.svg)](https://github.com/enckse/mayhem/actions/workflows/build.yml)
