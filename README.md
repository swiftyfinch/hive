# Hive 🐝

It's the prototype of the system to verify dependencies level.<br>
I used this as a playground during my self-education in the Go language.

Right now, it can parse only the Podfile.lock file.<br>
There are two subcommands: `tidy` and `check`.

## Module Types

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

## Dependency Rules

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

## Tidy

Call `tidy` to collect modules and dependencies in the first time.<br>
```sh
> hive tidy
```

It will create a config file which you should fill with module types. For example:
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

## Check

The `check` subcommand verifies local module dependencies based on default rules.
```sh
> hive check
⛔️ [api: base] LocalPodIO → Kingfisher
```
