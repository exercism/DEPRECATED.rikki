# Rikki

Rikki is the friendly neighborhood robot nitpicker.

## Usage

```bash
$ go build ./...
$ ./rikki
```

By default the worker runs against

- **redis**: localhost:6379
- **exercism**: localhost:4567
- **analysseur**: localhost:8989

These can be overridden using command-line flags:

```bash
$ ./rikki -redis=redis://user:pass@host:port/db/ -exercism=http://exercism.io -analysseur=http://analysseur.exercism.io
```
