# gh-azure-oidc
This is a Github CLI extension to setup a federated OIDC based connection between a repository, organization or an environment and Azure. 

## Usage instructions
```
gh azure-oidc [flags]

Options
--o <organization>
Select a organization to connect to azure
--R <[HOST/]OWNER/REPO>
Select a repository to connect to azure using the OWNER/REPO format
--e <environment>
Select a environment under the repository 
--useDefaults
This will skip the interactive flow and select the defaults on azure side to setup connection quickly (Not implemented yet, work in progress)


Examples
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


