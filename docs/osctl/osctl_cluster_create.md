<!-- markdownlint-disable -->
## osctl cluster create

Creates a local docker-based kubernetes cluster

### Synopsis

Creates a local docker-based kubernetes cluster

```
osctl cluster create [flags]
```

### Options

```
      --cpus string                 the share of CPUs as fraction (each container) (default "1.5")
  -h, --help                        help for create
      --image string                the image to use (default "docker.io/autonomy/talos:v0.3.0-alpha.7-9-g1e9a09e6-dirty")
      --kubernetes-version string   desired kubernetes version to run (default "1.16.2")
      --masters int                 the number of masters to create (default 1)
      --memory int                  the limit on memory usage in MB (each container) (default 1024)
      --mtu string                  MTU of the docker bridge network (default "1500")
      --workers int                 the number of workers to create (default 1)
```

### Options inherited from parent commands

```
      --name string          the name of the cluster (default "talos_default")
      --talosconfig string   The path to the Talos configuration file (default "/root/.talos/config")
  -t, --target strings       target the specificed node
```

### SEE ALSO

* [osctl cluster](osctl_cluster.md)	 - A collection of commands for managing local docker-based clusters

###### Auto generated by spf13/cobra on 13-Nov-2019
