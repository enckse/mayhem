Mayhem
===

A fork of [mayhem](https://github.com/BOTbkcd/mayhem) with some cut features
and generally simplified layout.

- No recurring tasks
- No steps
- JSON backing storage (use with other tools, hand edit, etc.)
- Standardized key bindings
- Cleaned up UI

[![build](https://github.com/enckse/mayhem/actions/workflows/build.yml/badge.svg)](https://github.com/enckse/mayhem/actions/workflows/build.yml)

## usage

### configuration

`mayhem` is configured via TOML (in `MAYHEM_CONFIG`, `XDG_CONFIG_HOME`, or `$HOME/.config/mayhem/`) in a settings.toml file

```
[data]
# override the location where data is stored
directory="~/.mayhem"
# save the data in a pretty (e.g. JSON pretty) indented/format
pretty=true

[display]
# display finished tasks that have been updated since
finished.since= "48h"

[backups]
# enable backups into a directory (offset from data.directory)
# backups are taken when mayhem starts
directory="backups"
# keeps this duration of backups
duration="72h"
# control the format of the date, allows controlling how many backups one gets
format="20060102"
```

### usage

Run `mayhem` and follow the navigation keys/help

## build

clone and `make`

[![build](https://github.com/enckse/mayhem/actions/workflows/build.yml/badge.svg)](https://github.com/enckse/mayhem/actions/workflows/build.yml)
