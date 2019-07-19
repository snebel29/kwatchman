### Testing
```shell
$ make test
```

### Dependencies
Will ensure dependencies using [dep](https://github.com/golang/dep), be careful with non locked dependencies
```shell
$ make deps
```
### Build
```shell
$ VERSION=1.0.0 make build
```
#### Docker image
```shell
$ VERSION=1.0.0 make docker-image
```
### Release
Relase of new docker images is achieved by creating a git tag, and through circleci
```shell
$ git tag v1.0.0
$ git push origin master --tags
```
