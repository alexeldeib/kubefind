# kubefind

A demo project to simply select and list arbitrary resource types by selector.

## Usage

`./kubefind --gvk "apps.v1.Deployment" --label tier=node`

will print all deployments matching the label `tier=node`.

gvk and label are both array flags.

You can specify several labels if you want to match a more granular set.

The output will be all objects from all specified GVKs matching the union of input labels.
