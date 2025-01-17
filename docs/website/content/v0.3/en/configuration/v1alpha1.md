---
title: v1alpha1
---

<!-- markdownlint-disable MD024 -->

Package v1alpha1 configuration file contains all the options available for configuring a machine.

We can generate the files using `osctl`.
This configuration is enough to get started in most cases, however it can be customized as needed.

```bash
osctl config generate --version v1alpha1 <cluster name> <cluster endpoint>
````

This will generate a machine config for each node type, and a talosconfig.
The following is an example of an `init.yaml`:

```yaml
version: v1alpha1
machine:
  type: init
  token: 5dt69c.npg6duv71zwqhzbg
  ca:
    crt: <base64 encoded Ed25519 certificate>
    key: <base64 encoded Ed25519 key>
  certSANs: []
  kubelet: {}
  network: {}
  install:
    disk: /dev/sda
    image: docker.io/autonomy/installer:latest
    bootloader: true
    wipe: false
    force: false
cluster:
  controlPlane:
    version: 1.16.2
    endpoint: https://1.2.3.4
  clusterName: example
  network:
    cni: ""
    dnsDomain: cluster.local
    podSubnets:
    - 10.244.0.0/16
    serviceSubnets:
    - 10.96.0.0/12
  token: wlzjyw.bei2zfylhs2by0wd
  certificateKey: 20d9aafb46d6db4c0958db5b3fc481c8c14fc9b1abd8ac43194f4246b77131be
  aescbcEncryptionSecret: z01mye6j16bspJYtTB/5SFX8j7Ph4JXxM2Xuu4vsBPM=
  ca:
    crt: <base64 encoded RSA certificate>
    key: <base64 encoded RSA key>
  apiServer: {}
  controllerManager: {}
  scheduler: {}
  etcd:
    ca:
      crt: <base64 encoded RSA certificate>
      key: <base64 encoded RSA key>
```

### Config

#### version

Indicates the schema used to decode the contents.

Type: `string`

Valid Values:

- ``v1alpha1``

#### machine

Provides machine specific configuration options.

Type: `MachineConfig`

#### cluster

Provides cluster specific configuration options.

Type: `ClusterConfig`

---

### MachineConfig

#### type

Defines the role of the machine within the cluster.

##### Init

Init node type designates the first control plane node to come up.
You can think of it like a bootstrap node.
This node will perform the initial steps to bootstrap the cluster -- generation of TLS assets, starting of the control plane, etc.

##### Control Plane

Control Plane node type designates the node as a control plane member.
This means it will host etcd along with the Kubernetes master components such as API Server, Controller Manager, Scheduler.

##### Worker

Worker node type designates the node as a worker node.
This means it will be an available compute node for scheduling workloads.

Type: `string`

Valid Values:

- ``init``
- ``controlplane``
- ``join``

#### token

The `token` is used by a machine to join the PKI of the cluster.
Using this token, a machine will create a certificate signing request (CSR), and request a certificate that will be used as its' identity.

Type: `string`

Examples:

```yaml
token: 328hom.uqjzh6jnn2eie9oi
```

> Warning: It is important to ensure that this token is correct since a machine's certificate has a short TTL by default

#### ca

The root certificate authority of the PKI.
It is composed of a base64 encoded `crt` and `key`.

Type: `PEMEncodedCertificateAndKey`

Examples:

```yaml
ca:
  crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJIekNCMHF...
  key: LS0tLS1CRUdJTiBFRDI1NTE5IFBSSVZBVEUgS0VZLS0tLS0KTUM...

```

#### certSANs

Extra certificate subject alternative names for the machine's certificate.
By default, all non-loopback interface IPs are automatically added to the certificate's SANs.

Type: `array`

Examples:

```yaml
certSANs:
  - 10.0.0.10
  - 172.16.0.10
  - 192.168.0.10

```

#### kubelet

Used to provide additional options to the kubelet.

Type: `KubeletConfig`

Examples:

```yaml
kubelet:
  image:
  extraArgs:
    key: value

```

#### network

Used to configure the machine's network.

Type: `NetworkConfig`

Examples:

```yaml
network:
  hostname: worker-1
  interfaces:
  nameservers:
    - 9.8.7.6
    - 8.7.6.5

```

#### disks

Used to partition, format and mount additional disks.
Since the rootfs is read only with the exception of `/var`, mounts are only valid if they are under `/var`.
Note that the partitioning and formating is done only once, if and only if no existing  partitions are found.

Type: `array`

Examples:

```yaml
disks:
  - device: /dev/sdb
    partitions:
      - size: 10000000000
        mountpoint: /var/lib/extra

```

> Note: `size` is in units of bytes.

#### install

Used to provide instructions for bare-metal installations.

Type: `InstallConfig`

Examples:

```yaml
install:
  disk:
  extraDiskArgs:
  extraKernelArgs:
  image:
  bootloader:
  wipe:
  force:

```

#### files

Allows the addition of user specified files.
Note that the file contents are not required to be base64 encoded.

Type: `array`

Examples:

```yaml
kubelet:
  contents: |
    ...
  permissions: 0666
  path: /tmp/file.txt

```

> Note: The specified `path` is relative to `/var`.

#### env

The `env` field allows for the addition of environment variables to a machine.
All environment variables are set on the machine in addition to every service.

Type: `Env`

Valid Values:

- ``GRPC_GO_LOG_VERBOSITY_LEVEL``
- ``GRPC_GO_LOG_SEVERITY_LEVEL``
- ``http_proxy``
- ``https_proxy``
- ``no_proxy``

Examples:

```yaml
env:
  GRPC_GO_LOG_VERBOSITY_LEVEL: "99"
  GRPC_GO_LOG_SEVERITY_LEVEL: info
  https_proxy: http://SERVER:PORT/

```

```yaml
env:
  GRPC_GO_LOG_SEVERITY_LEVEL: error
  https_proxy: https://USERNAME:PASSWORD@SERVER:PORT/

```

```yaml
env:
  https_proxy: http://DOMAIN\\USERNAME:PASSWORD@SERVER:PORT/

```

#### time

Used to configure the machine's time settings.

Type: `TimeConfig`

Examples:

```yaml
time:
  servers:
    - time.cloudflare.com

```

---

### ClusterConfig

#### controlPlane

Provides control plane specific configuration options.

Type: `ControlPlaneConfig`

Examples:

```yaml
controlPlane:
  version: 1.16.2
  endpoint: https://1.2.3.4
  localAPIServerPort: 443

```

#### clusterName

Configures the cluster's name.

Type: `string`

#### network

Provides cluster network configuration.

Type: `ClusterNetworkConfig`

Examples:

```yaml
network:
  cni: flannel
  dnsDomain: cluster.local
  podSubnets:
  - 10.244.0.0/16
  serviceSubnets:
  - 10.96.0.0/12

```

#### token

The [bootstrap token](https://kubernetes.io/docs/reference/access-authn-authz/bootstrap-tokens/).

Type: `string`

Examples:

```yaml
wlzjyw.bei2zfylhs2by0wd
```

#### aescbcEncryptionSecret

The key used for the [encryption of secret data at rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/).

Type: `string`

Examples:

```yaml
z01mye6j16bspJYtTB/5SFX8j7Ph4JXxM2Xuu4vsBPM=
```

#### ca

The base64 encoded root certificate authority used by Kubernetes.

Type: `PEMEncodedCertificateAndKey`

Examples:

```yaml
ca:
  crt: LS0tLS1CRUdJTiBDRV...
  key: LS0tLS1CRUdJTiBSU0...

```

#### apiServer

API server specific configuration options.

Type: `APIServerConfig`

Examples:

```yaml
apiServer:
  image: ...
  extraArgs:
    key: value
  certSANs:
    - 1.2.3.4
    - 5.6.7.8

```

#### controllerManager

Controller manager server specific configuration options.

Type: `ControllerManagerConfig`

Examples:

```yaml
controllerManager:
  image: ...
  extraArgs:
    key: value

```

#### scheduler

Scheduler server specific configuration options.

Type: `SchedulerConfig`

Examples:

```yaml
scheduler:
  image: ...
  extraArgs:
    key: value

```

#### etcd

Etcd specific configuration options.

Type: `EtcdConfig`

Examples:

```yaml
etcd:
  ca:
    crt: LS0tLS1CRUdJTiBDRV...
    key: LS0tLS1CRUdJTiBSU0...
  image: ...

```

---

### KubeletConfig

#### image

The `image` field is an optional reference to an alternative hyperkube image.

Type: `string`

Examples:

```yaml
image: docker.io/<org>/hyperkube:latest
```

#### extraArgs

The `extraArgs` field is used to provide additional flags to the kubelet.

Type: `map`

Examples:

```yaml
extraArgs:
  key: value

```

---

### NetworkConfig

#### hostname

Used to statically set the hostname for the host.

Type: `string`

#### interfaces

`interfaces` is used to define the network interface configuration.
By default all network interfaces will attempt a DHCP discovery.
This can be further tuned through this configuration parameter.

##### machine.network.interfaces.interface

This is the interface name that should be configured.

##### machine.network.interfaces.cidr

`cidr` is used to specify a static IP address to the interface.
This should be in proper CIDR notation ( `192.168.2.5/24` ).

> Note: This option is mutually exclusive with DHCP.

##### machine.network.interfaces.dhcp

`dhcp` is used to specify that this device should be configured via DHCP.

The following DHCP options are supported:

- `OptionClasslessStaticRoute`
- `OptionDomainNameServer`
- `OptionDNSDomainSearchList`
- `OptionHostName`

> Note: This option is mutually exclusive with CIDR.

##### machine.network.interfaces.ignore

`ignore` is used to exclude a specific interface from configuration.
This parameter is optional.

##### machine.network.interfaces.routes

`routes` is used to specify static routes that may be necessary.
This parameter is optional.

Routes can be repeated and includes a `Network` and `Gateway` field.

Type: `array`

#### nameservers

Used to statically set the nameservers for the host.
Defaults to `1.1.1.1` and `8.8.8.8`

Type: `array`

---

### InstallConfig

#### disk

The disk used to install the bootloader, and ephemeral partitions.

Type: `string`

Examples:

```yaml
/dev/sda
```

```yaml
/dev/nvme0
```

#### extraKernelArgs

Allows for supplying extra kernel args to the bootloader config.

Type: `array`

Examples:

```yaml
extraKernelArgs:
  - a=b

```

#### image

Allows for supplying the image used to perform the installation.

Type: `string`

Examples:

```yaml
image: docker.io/<org>/installer:latest

```

#### bootloader

Indicates if a bootloader should be installed.

Type: `bool`

Valid Values:

- `true`
- `yes`
- `false`
- `no`

#### wipe

Indicates if zeroes should be written to the `disk` before performing and installation.
Defaults to `true`.

Type: `bool`

Valid Values:

- `true`
- `yes`
- `false`
- `no`

#### force

Indicates if filesystems should be forcefully created.

Type: `bool`

Valid Values:

- `true`
- `yes`
- `false`
- `no`

---

### TimeConfig

#### servers

Specifies time (ntp) servers to use for setting system time.
Defaults to `pool.ntp.org`

> Note: This parameter only supports a single time server

Type: `array`

---

### Endpoint

---

### ControlPlaneConfig

#### version

Indicates which version of Kubernetes for all control plane components.

Type: `string`

Examples:

```yaml
1.16.2
```

> Note: The version must be of the format `major.minor.patch`, _without_ a leading `v`.

#### endpoint

Endpoint is the canonical controlplane endpoint, which can be an IP address or a DNS hostname.
It is single-valued, and may optionally include a port number.

Type: `Endpoint`

Examples:

```yaml
https://1.2.3.4:443
```

#### localAPIServerPort

The port that the API server listens on internally.
This may be different than the port portion listed in the endpoint field above.
The default is 6443.

Type: `int`

---

### APIServerConfig

#### image

The container image used in the API server manifest.

Type: `string`

#### extraArgs

Extra arguments to supply to the API server.

Type: `map`

#### certSANs

Extra certificate subject alternative names for the API server's certificate.

Type: `array`

---

### ControllerManagerConfig

#### image

The container image used in the controller manager manifest.

Type: `string`

#### extraArgs

Extra arguments to supply to the controller manager.

Type: `map`

---

### SchedulerConfig

#### image

The container image used in the scheduler manifest.

Type: `string`

#### extraArgs

Extra arguments to supply to the scheduler.

Type: `map`

---

### EtcdConfig

#### image

The container image used to create the etcd service.

Type: `string`

#### ca

The `ca` is the root certificate authority of the PKI.
It is composed of a base64 encoded `crt` and `key`.

Type: `PEMEncodedCertificateAndKey`

Examples:

```yaml
ca:
  crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJIekNCMHF...
  key: LS0tLS1CRUdJTiBFRDI1NTE5IFBSSVZBVEUgS0VZLS0tLS0KTUM...

```

#### extraArgs

Extra arguments to supply to etcd.
Note that the following args are blacklisted:

- `name`
- `data-dir`
- `initial-cluster-state`
- `listen-peer-urls`
- `listen-client-urls`
- `cert-file`
- `key-file`
- `trusted-ca-file`
- `peer-client-cert-auth`
- `peer-cert-file`
- `peer-trusted-ca-file`
- `peer-key-file`

Type: `map`

Examples:

```yaml
extraArgs:
  initial-cluster: https://1.2.3.4:2380
  advertise-client-urls: https://1.2.3.4:2379

```

---

### ClusterNetworkConfig

#### cni

The CNI used.

Type: `string`

Valid Values:

- `flannel`

#### dnsDomain

The domain used by Kubernetes DNS.
The default is `cluster.local`

Type: `string`

Examples:

```yaml
cluser.local
```

#### podSubnets

The pod subnet CIDR.

Type: `array`

Examples:

```yaml
podSubnets:
  - 10.244.0.0/16

```

#### serviceSubnets

The service subnet CIDR.

Type: `array`

Examples:

```yaml
serviceSubnets:
  - 10.96.0.0/12

```

---

### Bond

#### mode

The bond mode.

Type: `string`

#### hashpolicy

The hash policy.

Type: `string`

#### lacprate

The LACP rate.

Type: `string`

#### interfaces

The interfaces if which the bond should be comprised of.

Type: `array`

---

### Route

#### network

TODO.

Type: `string`

#### gateway

TODO.

Type: `string`

---
