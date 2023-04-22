# Application OpenAI plugin ✨

This application works as a generator tool that uses OpenAI to generate code
files. It takes a prompt as input and keeps generating code file content based
on the prompt until the user applies the generated files. if the content of the
files don't fully reflect what you expected you can keep adding context to
OpenAI until you get the desired result.

<!-- ## Installation

- Download the binary from [GitHub releases](https://github.com/afrancoc2000/application-ai/releases). -->

## Usage

### Prerequisites

`kubectl-ai` requires an [OpenAI API key](https://platform.openai.com/overview)
or an [Azure OpenAI Service](https://aka.ms/azure-openai) API key and endpoint.
Also, `kubectl` is required with a valid kubeconfig.

For both OpenAI and Azure OpenAI, you can use the following environment
variables:

```shell
export OPENAI_API_KEY=<your OpenAI key>
export OPENAI_DEPLOYMENT_NAME=<your OpenAI deployment/model name. defaults to "gpt-3.5-turbo">
```

> Following models are supported:
>
> - `code-davinci-002`
> - `text-davinci-003`
> - `gpt-3.5-turbo-0301` (deployment must be named `gpt-35-turbo-0301` for Azure
>   )
> - `gpt-3.5-turbo`
> - `gpt-35-turbo-0301`
> - `gpt-4-0314`
> - `gpt-4-32k-0314`

For Azure OpenAI Service, you can use the following environment variables:

```shell
export AZURE_OPENAI_ENDPOINT=<your Azure OpenAI endpoint, like "https://my-aoi-endpoint.openai.azure.com">
```

If `AZURE_OPENAI_ENDPOINT` variable is set, then it will use the Azure OpenAI
Service. Otherwise, it will use OpenAI API.

### Flags and environment variables

- `--skip-confirmation` flag or `SKIP_CONFIRMATION` environment variable can be
  set to prompt the user for confirmation before applying the manifest. Defaults
  to false.

- `--temperature` flag or `TEMPERATURE` environment variable can be set between
  0 and 1. Higher temperature will result in more creative completions. Lower
  temperature will result in more deterministic completions. Defaults to 0.

- `--chatContext` flag or `CHAT_CONTEXT` environment variable can be set between
  to add more context to the query. Defaults to "".

### How to use it

To use this tool, you need to run the `application-ai` app with a prompt as
argument. This prompt will be used to generate code files. The tool will keep
generating code file content based on the prompt until the user applies and the
files are generated. If the user decides not to apply the generated files, the
tool will exit without creating any files.

## Examples

Here is an example of how to use this tool:

```shell
$ "Create a project to deploy an AKS"
These are the files that would be created. Do you want to apply them? or add something to the query?
1. File: ./main.tf:

# Configure the Azure provider
provider "azurerm" {
        features {}
}

# Create a resource group
resource "azurerm_resource_group" "aks" {
        name     = var.resource_group_name
        location = var.resource_group_location
}

# Create an AKS cluster
resource "azurerm_kubernetes_cluster" "aks" {
        name                = var.cluster_name
        location            = azurerm_resource_group.aks.location
        resource_group_name = azurerm_resource_group.aks.name

        default_node_pool {
                name       = "default"
                vm_size    = var.node_vm_size
                node_count = var.node_count
        }

        identity {
                type = "SystemAssigned"
        }

        addon_profile {
                kube_dashboard {
                        enabled = true
                }
        }

output "password" {
        value = azurerm_kubernetes_cluster.aks.kube_config.0.password
}

output "username" {
        value = azurerm_kubernetes_cluster.aks.kube_config.0.username
}


Use the arrow keys to navigate: ↓ ↑ → ←
? Would you like to apply this? [Add to the query/Apply/Don't apply]:
+   Add to the query
  > Apply
    Don't apply
```

You can see that the application is suggesting to create a file named `main.tf`
with the code required to deploy an AKS in Azure, if you select
`Add to the query` you can specify further what the file should include and how
the app should behave, if you select `Apply` the file will be created, if you
choose `Don't apply` nothing will be created.

## Acknowledgements and Credits

Thanks to @sozercan for their work on Azure OpenAI fork in
https://github.com/sozercan/kubectl-ai which is based on https://github.com/simongottschlag/azure-openai-gpt-slack-bot which is based on https://github.com/PullRequestInc/go-gpt3
