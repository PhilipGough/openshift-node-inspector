OpenShift Node Inspector
=========================

CLI tool written in Go for fast debugging of Node applications running on OpenShift

### How it works
It works by mounting the Node Inspector source code along with its dependencies inside a container as a Volume.
It then starts a Node Inspector server and overwrites the CMD in the components Dockerfile to start the application with the debug flag.

A snapshot is taken of the current deployment configuaration and service which is then edited to update these objects. A route is created to allow
the Node Inspector to be exposed for that service.

### Prerequisites
1. [OpenShift CLI tool is installed](https://docs.openshift.com/enterprise/3.2/cli_reference/get_started_cli.html#installing-the-cli)
2. The user is logged and is using the appropriate project
3. The appropriate Security Context Constraint has "gitRepo" allowed as a Volume

### Install from Binary

Binaries are provided for Linux and Mac OS

_For Linux_

```bash
wget https://github.com/PhilipGough/openshift-node-inspector/releases/download/v0.1.0-alpha/oni-linux-amd64 -O /usr/bin/openshift-node-inspector

cd /usr/bin && sudo chmod 0777 openshift-node-inspector
```

_For Mac_

```bash
wget https://github.com/PhilipGough/openshift-node-inspector/releases/download/v0.1.0-alpha/oni-darwin-amd64 -O /usr/bin/openshift-node-inspector

cd /usr/bin && sudo chmod 0777 openshift-node-inspector
```

### Building

```bash
go get github.com/PhilipGough/openshift-node-inspector
````

### Using the tool

The CLI tool takes two basic commands with various options those are ```debug``` and ```clean``` followed the name of the component
we want to debug and some optional arguements as flags

#### The debug command

```bash
openshift-node-inspector [component-name] [options]
```

| Option   | Alias | Required | Default                                             | Description                                          |
|----------|-------|----------|-----------------------------------------------------|------------------------------------------------------|
| --port   | -p    | false    | 9000                                                | Port to listen on for Node Inspector's web interface |
| --image  | -i    | false    | Value taken from current deployment configuration   | Docker image to use                                  |
| --src    | -s    | false    | github.com/PhilipGough/openshift-node-inspector-src | Source code to mount as volume                       |
| --commit | -c    | false    | master                                              | Commit hash from Git repository                      |

Example

``` bash
openshift-node-inspector debug fh-supercore -p 8999
```

You will be prompted to accept the image or proved a different one

#### The clean command

The clean command simply attempts to revert to your previous state before debugging. It takes no options and must only proved the component

Example to revert the above debug

``` bash
openshift-node-inspector clean fh-supercore
```
