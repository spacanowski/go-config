# Config loader

Config loader let's you load config from yaml file into your apps config structs based on profile.

## How it works

Config loader takes `application` yaml file (either `yaml` or `yml` extension) and maps it into coresponding struct. You could overwrite or define additional params in `application-<profile>.yml` files where `<profile>` is any profile name, e.g. `application-dev.yml`, `application-test.yml` or any other. Profile flag and files are optional. Profile name arg can be passed by program call param `--profile=<profile>` or `-p=<profile>` where `<profile>` is profile name coresponding to one defined in profile config file.

If struct config value is not deffined in yaml file, then default type value is used. If yaml file has values not deffied in struct, they'll be ignored.

## Example usage

See [full example]() or short one below. // TODO add repo with example

Define config struct

```golang
type myAppConfig struct {
    appName      string
    accountName  string
    aws          struct {
        accountNumber int
        regions       []string
        clientSqs     map[string]string
        db            struct {
            autoCommit bool
            url        string
            username   string
        }
    }
}
```

Create `application.yml` file with common config

```yaml
appName: 'My app'
accountName: dev-account
aws:
  regions: ['eu-west-1']
  accountNumber: 123123
  clientSqs:
    client1: 'client1-q'
    client2: 'client2-q'
  db:
    autoCommit: true
    url: 'dev.db.url'
    username: db-user
```

Create optional profile config file `application-<profile>.yml`, e.g.

`application-test.yaml`
```yaml
accountName: test-account
aws:
  regions: ['eu-west-1']
  accountNumber: 321321
  clientSqs:
    client1: 'client1-q'
    client2: 'client2-q'
  db:
    url: 'test.db.url'
    username: test-db-user
```

`application-prod.yml`
```yml
accountName: prod-account
aws:
  regions: ['eu-west-1', 'eu-central-1']
  accountNumber: 441441
  clientSqs:
    client1: 'client1-q'
    client2: 'client2-q'
    client3: 'client3-q'
  db:
    url: 'prod.db.url'
    username: prod-db-user
```

Import loader and load config
// TODO add import

```go
func main() {
    var config myAppConfig
    if err := configloader.Load(&config); err != nil {
        // Handle config load error
    }
    // Use config
}
```

Run your app without profile (values taken only from `application.yml`)

`go run my-app.go`

or with profile app arg, depending on which config profile you want to use (values taken from `application.yml` and coresponding profile).

`go run my-app.go --profile=test`
or
`./my-app -p=test`
