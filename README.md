# Funplane - Fun Control Plane

Funplane provides a simple way to control envoy's sidecar configuration. It makes possible to tweak low level envoy
configuration and use xDS interface.


## Reference
Funplane uses kubernetes CustomResourceDefinitions to ship
envoy configuration into sidecars.

![Funplane Flow](img/Funplane.png)


![Controller Example](https://github.com/kubernetes/sample-controller/raw/master/docs/images/client-go-controller-interaction.jpeg)

## Quick start
```
$ funplane -c mesh.yaml -p 18000 -l info
```

## Development
Fanplane uses kubernetes crds toolset to generate clients and crd definition
For update the model use:

```
$ ./hack/update_codegen.sh
```
To generate go-client and deepcopy

## Roadmap

- Create registry
- Implement registry fetch
- Implement stream
- Improve tracing
- Implement tls
