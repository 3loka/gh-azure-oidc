/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/3loka/gh-azure-oidc/azureapi"
	"github.com/cli/go-gh"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cfgFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type setOptions struct {
	orgrepo string
}

func init() {

	var flag string
	rootCmd.Flags().StringVarP(&flag, "useDefaults", "f", "noDefault", "Use Defaults to create a connection quickly")
	rootCmd.Flags().Lookup("useDefaults").NoOptDefVal = "yes"
	rootCmd.PersistentFlags().String("o", "", "Select a organization to connect to azure")
	rootCmd.PersistentFlags().String("R", "", "Select a repository to connect to azure using the OWNER/REPO format")
	rootCmd.PersistentFlags().String("e", "", "Select a environment under the repository")
	rootCmd.PersistentFlags().String("useDefaults", "", "Use Defaults to create a connection quickly")

}

var rootCmd = &cobra.Command{
	Use:   "gh azure-oic",
	Short: "Connect Github to Azure for Workflow automation",
	Long:  `Connect Github to Azure for Workflow automation`,
	Run: func(cmd *cobra.Command, args []string) {
		orgrepo, _ := cmd.Flags().GetString("R")
		env, _ := cmd.Flags().GetString("e")
		orgFlag, _ := cmd.Flags().GetString("o")
		useDefaults, _ := cmd.Flags().GetString("useDefaults")
		if useDefaults == "yes" {
			fmt.Println("Use Defaults option is still work in progress, we are progressing with the non default flow for now")
		}
		runSetup(orgrepo, env, orgFlag, useDefaults)

	},
}

func runSetup(orgrepo string, env string, orgFlag string, useDefaults string) {

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	var repoName = orgrepo
	if repoName == "" {
		response := struct{ Login string }{}
		err = client.Get("user", &response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("running as %s\n", response.Login)

		repoNameContent := promptContent{
			fmt.Sprintf("Enter the org/repo? "),
			fmt.Sprintf("Enter the org/repo? "),
		}

		repoName = promptGetInput(repoNameContent)
	}

	split := strings.Split(repoName, "/")

	org := split[0]
	repo := split[1]

	// repo := flag.String("r", "default", "some name")
	// fmt.Println(*repo)

	fmt.Println("You will be redirected to Browser for login. Please complete the login and come back to the terminal")
	time.Sleep(2 * time.Second)

	token := azureapi.AuthenticateWithImplicitFlow()

	// token := azureapi.Authenticate()

	fmt.Println("Getting all tenants")

	// var tenantId, tenantname, subsriptionId, subscriptionName, ResourceGroupId, ResourceGroupName string

	tenantMap := azureapi.GetAllTenantsMap(token.AccessToken)

	directoryContent := promptContent{
		fmt.Sprintf("Choose your Azure Directory? "),
		fmt.Sprintf("Choose your Azure Directory? "),
	}

	tenantArray := make([]string, 0, len(tenantMap))
	for ind, _ := range tenantMap {
		tenantArray = append(tenantArray, ind)
	}

	directoryName := promptGetSelect(directoryContent, tenantArray)
	fmt.Println("Selected Directory: " + directoryName)
	tenantId := tenantMap[directoryName]

	subMap := azureapi.GetAllSubscriptions(token.AccessToken)

	subcontent := promptContent{
		fmt.Sprintf("Choose your Azure Subscription? "),
		fmt.Sprintf("Choose your Azure Subscription? "),
	}

	subArr := make([]string, 0, len(subMap))
	for ind, _ := range subMap {
		subArr = append(subArr, ind)
	}

	subName := promptGetSelect(subcontent, subArr)
	fmt.Println("Selected Subscription: " + subName)

	key := subMap[subName]
	createRg := promptContent{
		fmt.Sprintf("Time to pick a resource group. Do you want create a new one?"),
		fmt.Sprintf("Time to pick a resource group. Do you want create a new one?"),
	}
	createRgDecision := promptGetSelect(createRg, []string{"No", "Yes"})

	var resourceGroupId = ""

	if createRgDecision == "No" {
		rgContent := promptContent{
			fmt.Sprintf("Choose your Azure Resource Group associated to the subscription? "),
			fmt.Sprintf("Choose your Azure Resource Group associated to the subscription? "),
		}
		rgMap := azureapi.GetResourceGroupsPerSubscription(token.AccessToken, key)

		rgArr := make([]string, 0, len(rgMap))
		for ind, _ := range rgMap {
			rgArr = append(rgArr, ind)
		}

		rgName := promptGetSelect(rgContent, rgArr)
		fmt.Println("Selected Resource Group: " + rgName)
		resourceGroupId = rgName
	} else {
		rName := org + "-" + repo
		rNameContent := promptContent{
			fmt.Sprintf("Creating the Resource Group with Name: %v. Do you want to change the name?", rName),
			fmt.Sprintf("Creating the Resource Group with Name: %v. Do you want to change the name?", rName),
		}

		rNameDecision := promptGetSelect(rNameContent, []string{"No", "Yes"})
		if rNameDecision == "Yes" {
			rNameContent := promptContent{
				fmt.Sprintf("Enter the Resource Group Name you want to create"),
				fmt.Sprintf("Enter the Resource Group Name you want to create"),
			}

			rName = promptGetInput(rNameContent)
		}
		fmt.Println()
		fmt.Println(">>Creating Resource Group " + rName)
		resp := azureapi.CreateResourceGroup(token.AccessToken, key, rName)
		resourceGroupId = resp.Name
	}

	//newToken := azureapi.GetTokenWithTenantScope(token.RefreshToken)
	fmt.Println()
	fmt.Println("You will be redirected to Browser to grant additional role for FIC creation. Please complete the login and come back to the terminal")
	time.Sleep(2 * time.Second)

	newToken := azureapi.AuthenticateWithTenant(tenantId)
	// newToken := token
	//Create Azure Resources
	fmt.Println()
	fmt.Println(">>Creating Azure Application")
	// time.Sleep(2 * time.Second)
	appResponse := azureapi.CreateAzureApplication(newToken.AccessToken, repo)

	//Create Service Principal
	fmt.Println()
	fmt.Println(">>Creating Service Principal")
	// time.Sleep(2 * time.Second)
	servicePrincipal := azureapi.CreateServicePrincipal(newToken.AccessToken, appResponse.AppId)

	//Create FIC
	fmt.Println()
	fmt.Println(">>Creating Federated Identity Credentials")
	// time.Sleep(2 * time.Second)
	azureapi.CreateFIC(newToken.AccessToken, appResponse.Id, repo)

	//Assign Role Definition
	fmt.Println()
	fmt.Println(">>Assigning Role Definition")
	// time.Sleep(2 * time.Second)
	azureapi.AssignRoleDefinition(token.AccessToken, servicePrincipal.Id, key, resourceGroupId)

	//Creating secrets
	fmt.Println()
	fmt.Printf(">>Creating Client Secrets in %s \n", orgrepo)
	fmt.Printf("AZURE_CLIENT_ID: %s \n", appResponse.Id)
	fmt.Printf("AZURE_TENANT_ID Client: %s \n", tenantId)
	fmt.Printf("AZURE_SUBSCRIPTION_ID: %s \n", key)

	createSecret("AZURE_CLIENT_ID", appResponse.Id, orgrepo, orgFlag, env)
	createSecret("AZURE_TENANT_ID", tenantId, orgrepo, orgFlag, env)
	createSecret("AZURE_SUBSCRIPTION_ID", key, orgrepo, orgFlag, env)

}

func createSecret(name string, value string, orgrepo string, orgFlag string, env string) {
	args := []string{"secret", "set", name, "--body", value}
	if env != "" {
		args = append(args, "--env", env)
	}
	if orgFlag != "" {
		args = append(args, "--org", orgFlag)
	} else {
		args = append(args, "-R", orgrepo)
	}
	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stdOut.String())
}

type promptContent struct {
	errorMsg string
	label    string
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Input: %s\n", result)

	return result
}

func promptGetSelect(pc promptContent, items []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: pc.label,
			Items: items,
			// AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Input: %s\n", result)

	return result
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
