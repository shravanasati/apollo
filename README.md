# apollo

A simple CLI tool that allows to have recurring infinite timers.

It can be used as tool to remind you to do tasks, like taking a eye break every few minutes.

### Installation

```
go install github.com/Shravan-1908/apollo@latest
```

### Usage

Just call `apollo` from the command line.

This is the apollo config file located at `~/.apollo`.

```json
{
    "timeouts": {
        "exercise_timeout": 2400,
        "eyes_timeout": 900,
        "water_timeout": 1800
    },
    "play_beep": true,
    "notify": true
}
```

Modify it according to your requirements. All timeout durations are in seconds.
