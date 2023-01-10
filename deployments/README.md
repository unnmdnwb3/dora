# Helm Chart

You find the helm chart for `dora` in this directory. It has a dependecy to [Bitnami's MongoDB chart](https://bitnami.com/stack/mongodb/helm).

## Install

To install the helm chart, be sure to have helm configured and to have a running cluster ready. To install the chart

```bash
helm install dora ./deployments/dora --namespace dora --create-namespace
```

## Testing

To test the helm chart without deploying it, just use:

```bash
helm install dora ./deployments/dora --namespace dora --create-namespace --dry-run --debug   
```

If you want to test a running helm chart, e.g. on a `kind` cluster, just use the installation guide above.
