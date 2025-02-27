# Custom Resource Definitions

In this directory, you would place the CRD YAML files for your operator.

The standard approach is to copy your CRD YAML files directly from:
`config/crd/bases/monitoring.monitoring.example.com_observabilitystacks.yaml`

Per Helm best practices, CRDs should be installed separately from the chart's templates, which is why they are placed in this special directory.

When installing the chart, you would use:
```
helm install kube-insight-operator ./kube-insight-operator --set crd.create=true
```

Or, for more control, install the CRDs manually before installing the chart:
```
kubectl apply -f config/crd/bases/monitoring.monitoring.example.com_observabilitystacks.yaml
helm install kube-insight-operator ./kube-insight-operator --set crd.create=false
```