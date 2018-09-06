# Helm-Bulk

This is a Helm plugin that loads or saves Helm releases from File to Cluster,
or Cluster to File, respectively.

## Installation

Clone the repo:

```
$ git clone git@github.com:ovotech/helm-bulk.git
```

Build the binary:

```
$ cd helm-bulk

$ go build
```

The plugin hasn't been open sourced yet, so the absolute path to the local checked out git repo (rather than the git https url) must be given to Helm in order to install it.

```
$ helm plugin install <abs_path_to_local_git_dir>
```

This will:

1. Find the `plugin.yaml` in the directory you give it
2. Executes the file specified in `hooks: install:` (useful for any pre-install
  operations)
3. Copies the file that's specified in the `command` value into the Helm plugin directory (defaults to `.helm/plugins/`). This is the file that'll be executed when you invoke the plugin with Helm, i.e. `helm bulk`.

The final String of the output, on a successful install, will be: `Installed plugin: bulk`.


You can also verify it's been installed using:

```
$helm plugin list                                     

NAME	VERSION	DESCRIPTION
bulk	0.0.8  	Load or Save Helm Releases from File to Cluster, or Cluster to File, respectively
```

Try invoking it:

```
$ helm bulk

Load or Save Releases from File to Cluster, or Cluster to File, respectively

Usage:
  helm-bulk [command]

Available Commands:
  help        Help about any command
  load        Load Releases from File to Cluster
  save        Save Releases from Cluster to File
  show        Show Releases currently stored in the file

Flags:
  -f, --filename string   Filename to use with a Load or Save command (default "helm-releases.txt")
  -h, --help              help for helm-bulk

Use "helm-bulk [command] --help" for more information about a command.
```

## Getting Started

`helm-bulk` will only ever use your current kubectl context, so make sure
you've switched to whatever Context/Cluster you want to use (e.g. `kubectl
  config use-context <context_name` or `gcloud container clusters....` to
  re-auth into your target Cluster).

If end-to-end testing, try following these commands through in order, otherwise
they can be run individually:

```
# Use Helm to list your Releases
$ helm ls

# Save deployed Helm Releases to file (defaults to "./helm-releases.txt")
$ helm bulk save

# Print out a list of Helm Releases currently stored in the File
# (defaults to "./helm-releases.txt")
$ helm bulk show

######################################################
# if e-2-e testing, simulate loss of Helm Releases in Cluster here
# e.g.:
# recreate cluster, delete all current Helm Releases in Cluster, OR switch
# kubectl context to a fresh Cluster
######################################################

# Load Helm Releases from File into your Cluster
$ helm bulk load

# Use Helm to list the Releases again
$ helm ls
```
