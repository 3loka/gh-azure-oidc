package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/3loka/gh-azure-oidc/azureapi"
	"github.com/cli/go-gh"
	"github.com/manifoldco/promptui"
)

func main() {

	// cmd.Execute()
	
	fmt.Println("hi world, this is the gh-azure-oidc extension!")
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	repoName := promptGetInput(repoNameContent)

	split := strings.Split(repoName, "/")

	org := split[0]
	repo := split[1]

	// repo := flag.String("r", "default", "some name")
	// fmt.Println(*repo)

	fmt.Println("You will be redirected to Browser for login. Please complete the login and come back to the terminal")
	time.Sleep(3 * time.Second)

	accessToken := azureapi.Authenticate()
	tenantMap := azureapi.GetAllTenantsMap(accessToken)

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

	subMap := azureapi.GetAllSubscriptions(accessToken)

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
	createRgDecision := promptGetSelect(createRg, []string{"Yes", "No"})

	if createRgDecision == "No" {
		rgContent := promptContent{
			fmt.Sprintf("Choose your Azure Resource Group associated to the subscription? "),
			fmt.Sprintf("Choose your Azure Resource Group associated to the subscription? "),
		}
		rgMap := azureapi.GetResourceGroupsPerSubscription(accessToken, key)

		rgArr := make([]string, 0, len(rgMap))
		for ind, _ := range rgMap {
			rgArr = append(rgArr, ind)
		}

		rgName := promptGetSelect(rgContent, rgArr)
		fmt.Println("Selected Resource Group: " + rgName)
	} else {
		rName := org + "-" + repo
		fmt.Println("Creating Resource Group " + rName)
		azureapi.CreateResourceGroup(accessToken, key, rName)
	}

	//Create Azure Resources
	fmt.Println("Creating Azure Application")
	time.Sleep(2 * time.Second)
	appId := azureapi.CreateAzureApplication(accessToken, repo)

	//Create SP
	fmt.Println("Creating Service Principal")
	time.Sleep(2 * time.Second)
	azureapi.CreateServicePrincipal(accessToken, appId)

	//Create FIC
	fmt.Println("Creating Federated Identity Credentials")
	time.Sleep(2 * time.Second)

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

	fmt.Printf("Input: %s\n", result)

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

	fmt.Printf("Input: %s\n", result)

	return result
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
