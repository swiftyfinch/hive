# Hive 🐝

It's a prototype of a system to verify dependencies level.<br>
I used this as a playground during my self-education in the Go language.

Right now, it can parse only the Podfile.lock file.<br>
There are two subcommands: `tidy` and `check`.

## Tidy

Call `tidy` to collect modules and dependencies in the first time.<br>
```sh
hive tidy
```

It will create a config file which you can tune for yourself. For example:
```yml
types: # Describe dependencies type here
  - base
  - feature
  - tests: .*Tests? # You can use regular expression to auto-detect type
  - io: IO
bans: # Ban some dependencies
  - feature: feature
  - base: feature
  - base: base
    severity: warning
  - tests: feature
    severity: warning
  - io: base
modules:
  remote:
    Alamofire: base # Choose a type for each module
    Kingfisher: base
    Moya/Core: base
    SnapKit: base
  local:
    LocalPod: feature
    LocalPod/Tests: tests
    LocalPodIO: io
```

And use this command after each change of modules/dependencies/config.<br>
It will remove old modules and add new ones.

## Check

The `check` subcommand verifies local module dependencies based on `bans` field from the config.
```sh
hive check
```
```sh
[warning] LocalPod/Tests(tests) → LocalPod(feature)
[error] LocalPodIO(io) → Kingfisher(base)
```
