# frakti

[![GoDoc](https://godoc.org/k8s.io/frakti?status.svg)](https://godoc.org/k8s.io/frakti) [![Build Status](https://travis-ci.org/kubernetes/frakti.svg?branch=master)](https://travis-ci.org/kubernetes/frakti)

Frakti enables hypervisor-agnostic container runtime for Kubernetes via 
[HyperContainer](http://hypercontainer.io/). It provides a kubelet runtime API
which will be consumed by kubelet.

## Build

```sh
mkdir -p $GOPATH/src/k8s.io
git clone https://github.com/kubernetes/frakti.git $GOPATH/src/k8s.io/frakti
cd $GOPATH/src/k8s.io/frakti
make && make install
```

## Start frakti

First start hyperd with gRPC endpoint `127.0.0.1:22318`:

```sh
$ grep gRPC /etc/hyper/config
gRPCHost=127.0.0.1:22318
```

Then start frakti:

```sh
frakti -v=3 --logtostderr --listen=127.0.0.1:10238 --hyper-endpoint=127.0.0.1:22318
```

## Start kubelet with frakti

```sh
kubelet --container-runtime-endpoint=127.0.0.1:10238 ...
```

## Links

- [HyperContainer](http://hypercontainer.io/)
- [HyperContainer src](https://github.com/hyperhq/hyperd)
