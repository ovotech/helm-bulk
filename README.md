# Helm-Bulk

This is a Helm plugin that loads or saves Helm releases from File to Cluster,
or Cluster to File, respectively.

## Installation

Install the latest version:

```shell
$ helm plugin install https://github.com/ovotech/helm-bulk
```

Install a specific version:

```shell
$ helm plugin install https://github.com/ovotech/helm-bulk --version 0.0.20
```

This will:

1. Find the `plugin.yaml` in the project root
2. Execute the file specified in `hooks: install:` (useful for any pre-install
  operations)
3. Copy the file that's specified in the `command` value into the Helm plugin
directory (defaults to `.helm/plugins/`). This is the file that'll be executed
when you invoke the plugin with Helm, i.e. `helm bulk`.

The final String of the output, on a successful install, will be:
`helm-bulk [version] is correctly installed.`.


You can also verify it's been installed using:

```
$ helm plugin list                                     

NAME	VERSION	DESCRIPTION
bulk	0.0.20  	Load or Save Helm Releases from File to Cluster, or Cluster to
File, respectively
```

Try invoking it:

```
$ helm bulk

Load or Save Releases from File to Cluster, or Cluster to File, respectively
...
...
```

## Getting Started

`helm-bulk` will only ever use your current kubectl context, so make sure
you've switched to whatever Context/Cluster you want to use (e.g. `kubectl
  config use-context <context_name>` or `gcloud container clusters....` to
  re-auth into your target Cluster).

If end-to-end testing, try following these commands through in order, otherwise
they can be run individually:

```
# Use Helm to list your Releases
$ helm ls

# Save deployed Helm Releases to archive (defaults to "./helm-releases.tar.gz")
$ helm bulk save

# Print out a list of Helm Releases currently stored in the archive
# (defaults to "./helm-releases.tar.gz")
$ helm bulk show

###############################################################################
# if e-2-e testing, simulate loss of Helm Releases in Cluster here
# e.g.:
# recreate cluster, delete all current Helm Releases in Cluster, OR switch
# kubectl context to a fresh Cluster
###############################################################################

# Load Helm Releases from File into your Cluster
$ helm bulk load

# Use Helm to list the Releases again
$ helm ls
```

## Idempotency

`helm bulk load` will attempt to get the Helm Releases in your Cluster to
reflect what you have in your File, but only for those Releases defined. It
won't touch any Releases you may already have in your Cluster.

If you already have one or more of the Releases in your File installed in your
Cluster, then things get a little more complicated than working with an 'empty'
Cluster. `helm-bulk` needs to work out whether to ignore, delete or upgrade the
existing Releases.

By default, `helm-bulk` will ignore the existing Releases. If you want it to
delete or upgrade, use the `-d` or `-u` flags respectively.

`helm-bulk` is designed to be used shortly after Cluster create (obviously post
  tiller install), in which case there won't be any existing Helm Releases.

## Save order

By default, Releases will be saved to file in the order they are returned by
`helm ls`, which by default in Helm is alphabetically.

It's possible to specify an overriding order by dropping a yaml file such as:

```
order:
  - first-release
  - second-release
```

..into the working directory, or into a custom directory and specifying this
directory in the `-c, --order-pref-config-dir` flag when the plugin is called.
The config file **must** be named `orderPref.yaml`.

Alternatively the order can be specified in an environment variable like so:

```
$ export HELM_BULK_ORDER=first-release,second-release
```

Note: the ordering is applied at time of `helm bulk save`. When `helm bulk load`
is called, the plugin will just blindly follow the order in the file from top
to bottom.

This could be useful in cases where all Releases are dependent on a small number
of other Releases, e.g. installing CRDs or application config.

Helm dependencies could also be used to achieve the same thing, and should in
theory work, although that hasn't been tested to any degree.

## Release Naming

When you install a Helm Chart, if you don't provide a name, Helm will generate
one for you, e.g. "kissing-wildebeest". Subsequent `helm install` commands
(also with no name specified) will create new Releases with different names
(even if for the same Chart). **All** of these Releases will then be returned
when running `helm ls`, and therefore will also make it into your File after a
`helm bulk save`.

This would lead to greater processing times when running a `helm bulk load` or
`helm bulk save`, and a larger persisted File, and could lead to unexpected
`load` results.

To prevent this, it's recommended to **always name your Release when
installing**, so you only have one current Release (and any upgrades to that
will supersede, with version number incrementing).

## Release History

History of Helm Releases is currently not persisted.
