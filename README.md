# kv demo project

This repo contains items for the key-value demo project.

Here is a short recording of the project in action:
[![asciicast](https://asciinema.org/a/Vy8UtezxttWwDk9nikVPFTjNy.png)](https://asciinema.org/a/Vy8UtezxttWwDk9nikVPFTjNy?speed=1.5&theme=solarized-light&size=small)


The project includes the following components:

* API
* CLI
* Infrastructure Components

## API
The API provides the mechanism for interacting with the key-value service. It
provides a small number of operations to write/read/delete data.

There are two domain objects exposed through the API:
* stores
* keys

Stores are a container for keys. A store can have `n` keys. Although not
implemented in this demo, stores could provide an entity upon which other
features could leverage (security, functionality, durability, etc.). When a
store is deleted, all of the keys contained within it are also deleted.

Keys are relatively straight forward. A key must belong to a single store. Keys
cannot be nested (i.e. `foo/bar/baz`). The value of a key is limited to 2K.

The API provides the following resources/methods for manipulating data:

* `PUT /stores/{store-name}` - create a store
* `DELETE /stores/{store-name}` - delete a store (and all keys within)
* `PUT /stores/{store-name}/keys/{key-name}` - set a key (both create and update)
* `GET /stores/{store-name}/keys/{key-name}` - get the value of a key
* `DELETE /stores/{store-name}/keys/{key-name}` - delete a key
* `GET /status` - basic endpoint used for load balancer health checks

The API uses [consul](https://consul.io) for it's persistence layer.

Have a look at the [Local Development
section](https://github.com/thbishop/kv#local-development) of this `README` to
find out how to work with the API locally.

## CLI
The CLI provides an easier way to interact with the key-value service API. It
provides commands for all functionality of the API. This includes:

* creating and deleting stores
* setting, getting, and deleting keys

There are pre-built binaries for [linux
(x64)](https://s3-us-west-2.amazonaws.com/kv-artifacts-us-west-2/kv-linux.zip) and
[osx/macOS](https://s3-us-west-2.amazonaws.com/kv-artifacts-us-west-2/kv-darwin.zip).

Once downloaded, you can run:

```sh
unzip kv-*.zip
chmod +x kv
./kv -h
```

## Infrastructure
The key-value service runs in AWS (us-west-2). The vast majority of the
infrastructure used to run the service was created using cfn. You can find the
cfn templates (along with a couple of install helper scripts) in the `./infra`
directory.

The cfn templates are mostly agnostic, but in the interest of saving time, there
are some pieces which are hardcoded which may prevent them from working in a
different account/region.

The key-value service API endpoint is at: `https://kv-api.dyson-sphere.com`

## Local Development
_NOTE: This repo assumes you'll be working in a `*nix` style environment._

### API
The API was built using golang 1.10.1 and uses
[dep](https://golang.github.io/dep/). There's a `Makefile` in the `./api`
directory which will help you build/test/package the API. After cloning this
repo into the proper `GOPATH` directory (`$GOPATH/src/github.com/thbishop/kv`),
you can run the following for help with the `Makefile`:

```sh
cd $GOPATH/src/github.com/thbishop/kv
cd api
make help
```

After building the API, you can start it with:
```sh
./bin/api
```

By default, the API will listen on port `8080`. If you want to change this, you
can set the `KV_API_PORT` environment variable:
```sh
env KV_API_PORT=8181 ./bin/api
```

As the API uses consul for persistence, it will connect to a consul agent
running on the local system using the default listenting port (8500). After
[downloading consul](https://www.consul.io/downloads.html), you can
start a local, single-node consul cluster with:

```sh
consul agent -server -bootstrap-expect=1 -data-dir=./data -node=agent-one
```

### CLI
The CLI was also built using golang 1.10.1 and uses
[dep](https://golang.github.io/dep/). There's a `Makefile` in the `./cli`
directory which will help you build/test/package the CLI. After cloning this
repo into the proper `GOPATH` directory (`$GOPATH/src/github.com/thbishop/kv`),
you can run the following for help with the `Makefile`:

```sh
cd $GOPATH/src/github.com/thbishop/kv
cd cli
make help
```

By default, the CLI uses the API endpoint in AWS. For local development, you can
override the API endpoint the CLI targets with:
```sh
env KV_CLI_API_URL=http://localhost:8080 bin/kv -h
```

### Basic Integration Test
This repo also contains a very basic integration test script. The script is
located in the `./test` directory. This integration test script leverages the
CLI and expects it to be available at `./cli/bin/kv` (which is where it is built
when using the CLI `Makefile`). You can run the script with:
```sh
./test/integration.sh
```

Or if you want to run it against a local API server, you can do so with:
```sh
# you may need to change the port if your API server is listening on a
# non-default port
env KV_CLI_API_URL=http://localhost:8080 ./test/integration.sh
```
