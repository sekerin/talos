---
title: 'Digital Ocean'
---

## Creating a Cluster via the CLI

In this guide we will create an HA Kubernetes cluster with 1 worker node.
We assume an existing [Space](https://www.digitalocean.com/docs/spaces/), and some familiarity with Digital Ocean.
If you need more information on Digital Ocean specifics, please see the [official Digital Ocean documentation](https://www.digitalocean.com/docs/).

### Create the Image

First, download the Digital Ocean image from a Talos release.

Using an upload method of your choice (`doctl` does not have Spaces support), upload the image to a space.
Now, create an image using the URL of the uploaded image:

```bash
doctl compute image create \
    --region $REGION \
    --image-name talos-digital-ocean-tutorial \
    --image-url https://talos-tutorial.$REGION.digitaloceanspaces.com/digital-ocean.raw.gz \
    Talos
```

Save the image ID.
We will need it when creating droplets.

### Create a Load Balancer

```bash
doctl compute load-balancer create \
    --region $REGION \
    --name talos-digital-ocean-tutorial-lb \
    --tag-name talos-digital-ocean-tutorial-control-plane \
    --health-check protocol:tcp,port:6443,check_interval_seconds:10,response_timeout_seconds:5,healthy_threshold:5,unhealthy_threshold:3 \
    --forwarding-rules entry_protocol:tcp,entry_port:443,target_protocol:tcp,target_port:6443
```

We will need the IP of the load balancer.
Using the ID of the load balancer, run:

```bash
doctl compute load-balancer get --format IP <load balancer ID>
```

Save it, as we will need it in the next step.

### Create the Machine Configuration Files

#### Generating Base Configurations

Using the DNS name of the loadbalancer created earlier, generate the base configuration files for the Talos machines:

```bash
$ osctl config generate talos-k8s-digital-ocean-tutorial https://<load balancer IP or DNS>
created init.yaml
created controlplane.yaml
created join.yaml
created talosconfig
```

At this point, you can modify the generated configs to your liking.

#### Validate the Configuration Files

```bash
$ osctl validate --config init.yaml --mode cloud
init.yaml is valid for cloud mode
$ osctl validate --config controlplane.yaml --mode cloud
controlplane.yaml is valid for cloud mode
$ osctl validate --config join.yaml --mode cloud
join.yaml is valid for cloud mode
```

### Create the Droplets

#### Create the Bootstrap Node

```bash
doctl compute droplet create \
    --region $REGION \
    --image <image ID> \
    --size s-2vcpu-4gb \
    --enable-private-networking \
    --tag-names talos-digital-ocean-tutorial-control-plane \
    --user-data-file init.yaml \
    --ssh-keys <ssh key fingerprint> \
    talos-control-plane-1
```

> Note: Although SSH is not used by Talos, Digital Ocean still requires that an SSH key be associated with the droplet.
> Create a dummy key that can be used to satisfy this requirement.

#### Create the Remaining Control Plane Nodes

Run the following twice, to give ourselves three total control plane nodes:

```bash
doctl compute droplet create \
    --region $REGION \
    --image <image ID> \
    --size s-2vcpu-4gb \
    --enable-private-networking \
    --tag-names talos-digital-ocean-tutorial-control-plane \
    --user-data-file controlplane.yaml \
    --ssh-keys <ssh key fingerprint> \
    talos-control-plane-2
doctl compute droplet create \
    --region $REGION \
    --image <image ID> \
    --size s-2vcpu-4gb \
    --enable-private-networking \
    --tag-names talos-digital-ocean-tutorial-control-plane \
    --user-data-file controlplane.yaml \
    --ssh-keys <ssh key fingerprint> \
    talos-control-plane-3
```

#### Create the Worker Nodes

Run the following to create a worker node:

```bash
doctl compute droplet create \
    --region $REGION \
    --image <image ID> \
    --size s-2vcpu-4gb \
    --enable-private-networking \
    --user-data-file join.yaml \
    --ssh-keys <ssh key fingerprint> \
    talos-worker-1
```

### Retrieve the `kubeconfig`

To configure `osctl` we will need the first controla plane node's IP:

```bash
doctl compute droplet get --format PublicIPv4 <droplet ID>
```

At this point we can retrieve the admin `kubeconfig` by running:

```bash
osctl --talosconfig talosconfig config target <control plane 1 IP>
osctl --talosconfig talosconfig kubeconfig > kubeconfig
```
