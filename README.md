# Hive 🐝

It's the prototype of the system to verify dependencies level.<br>
I used this as a playground during my self-education in the Go language.

Right now, it can parse only the Podfile.lock file.<br>


## Details

There is an internal declaration of module types and dependency rules.

### Module Types

```go
// Can use any regexp for different platforms
map[string][]string{
  "tests":   {".*Tests$"},
  "app":     {".*Example$"},
  "mock":    {".*Mock$"},
  "feature": {},
  "base":    {},
  "api":     {".*IO$", ".*Interfaces$"},
}
```

### Dependency Rules

```go
map[string][]string{
  "tests":   {"api", "base", "feature", "mock", "app"},
  "app":     {"api", "base", "feature", "mock"},
  "mock":    {"api", "base"},
  "feature": {"api", "base"},
  "base":    {"api", "base"},
  "api":     {},
}
```

```mermaid
flowchart LR
    tests --> app
    tests ---> feature
    tests ---> base
    tests ---> mock
    tests ---> api
    app ---> api
    app ---> feature
    app ---> base
    app ---> mock
    feature ---> base
    feature ---> api
    mock ---> base
    mock ---> api
    base --> base2[base]
```

## Commands

There are two subcommands: `tidy` and `check`.

### Tidy

Call `tidy` to collect modules and dependencies in the first time.<br>
And then call it to sync config with the current modules.
```sh
> hive tidy
```

It will create a config file with all modules. For example:
```yml
modules:
  remote:
    Alamofire: null
    Kingfisher: null
    Moya/Core: null
    SnapKit: null
  local:
    LocalPod: null
    LocalPod/Tests: tests
    LocalPodIO: api
    LocalPodMock: mock
    LocalPodsExample: app
```

And then fill missing types:
```yml
modules:
  remote:
    Alamofire: base
    Kingfisher: base
    Moya/Core: base
    SnapKit: base
  local:
    LocalPod: feature
    LocalPod/Tests: tests
    LocalPodIO: api
    LocalPodMock: mock
    LocalPodsExample: app
```

And use this command after each change of modules/dependencies/config.<br>
It will remove old modules and add new ones.

### Check

The `check` subcommand verifies local module dependencies based on default rules.
```sh
> hive check
⛔️ [api: base] LocalPodIO → Kingfisher
```
