# GCIDP Agent
## Note
This is a tool for internal use only! It is made public to easily test and integrate

## About
The goal of the project is to spin up containers for different branches/versions
of a project by simply defining a `.gcidp/build.go` file in the root directory of your project.

`build.go` is a regular Go file where all the syntax and standard library works as expected.
You can run whatever code you want, starting from simple `if` statements to complex pipelines.

The project also contains a set of utils to spin up docker containers, 
define build pipelines and many more features coming.


## Usage
Your `build.go` should contain the following code
```go
var BuildServerVersion = "0.2.2" // this tells GCIDP server which version of the agent you are using
// At the moment server checks if the major versions match (0.2.x == 0.2.x), otherwise the build fails.
// Future iterations will include a more robust versioning system

func Cleanup(runner *pipeline.Runner) {
	// this gets called when a branch is deleted
    // your code here
}

func Build(runner *pipeline.Runner) {
	// this gets called when a push to a branch is made
    // your code here
}
```

You can define pipelines, a set of stages executed consecutively
```go
runner.Pipeline(
	// your stages go here
)
```

If you want to executed stages concurrently, define multiple pipelines
```go
runner.Pipeline(
    // your stages go here
)
runner.Pipeline(
    // your stages go here
)
```

## Logging
You can log messages to the UI using the `Logger`
```go
runner.Logger.Debug("message")
runner.Logger.Info("message")
runner.Logger.Error("message")
```

## Docker
Building a docker image
```go
docker.Build(imageName, contextDir).Target("prod"),
```

Running a docker container
```go
docker.Run(containerName, imageName).Config(
    docker.PortBinding("5432", "5432"),
    docker.Volume("~/volumes/zender/projectName", "/var/lib/postgresql/data"),
    docker.Env("POSTGRES_DB", "postgres"),
    docker.Env("POSTGRES_PASSWORD", "postgres"),
    docker.Network("app"),
)
```

## Docker utils
Define environment variables
```go
docker.Env("KEY", "VALUE")
```
Equivalent to
```yaml
environment:
    -KEY=VALUE
```

Define port bindings
```go
docker.PortBinding("5432", "5432")
```

Equivalent to
```yaml
ports:
  - "5432:5432"
```

Define volumes
```go
docker.Volume("~/volumes/project/data", "/var/lib/postgresql/data")
```

Equivalent to
```yaml
volumes:
  - "~/volumes/project/data:/var/lib/postgresql/data"
```

Define networks
```go
docker.Network("app")
```

Equivalent to
```yaml
networks:
  - app
```

Set hostname
```go
docker.Hostname("postgres")
```

Equivalent to
```yaml
db:
  image: postgres
  hostname: postgres
```