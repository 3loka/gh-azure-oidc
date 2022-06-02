# gh-azure-oidc
This is a Github CLI extension to setup a federated OIDC based connection between a repository, organization or an environment and Azure resource groups. 

## Usage instructions
```
gh azure-oidc [flags]

    Options
        --o <organization>
        Select a organization to connect to azure
        --e <environment>
        Select an environment under the repository 
        --useDefaults
        This will skip the interactive flow and select the defaults on azure side to setup connection quickly (Not implemented yet, work in progress)

    Options inherited from parent commands
        -R, --repo <[HOST/]OWNER/REPO>
        Select another repository using the [HOST/]OWNER/REPO format


Examples
# Setup connection at a repo level (assuming user is in the git folder)
$ gh azure-oidc

# Setup connection at an organization level
$ gh azure-oidc -o myorg

# Setup connection at a repo level
$  gh azure-oidc -R myorg/myrepo

# Setup connection at a repo and environment level
$  gh azure-oidc -R myorg/myrepo -e myenvironment

```

![Demo](https://github.com/3loka/gh-azure-oidc/blob/main/azure-cli-demo.gif)

# To run this application, execute below command
- Install the CLI extension - Follow the steps here https://github.com/cli/cli#installation
- Clone this repo - `git clone https://github.com/3loka/gh-azure-oidc.git`
- Build the code - `go build`


