package azureapi

import (
	"fmt"
	"strings"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/3loka/gh-azure-oidc/auth"

	"github.com/3loka/gh-azure-oidc/models"
)

func Authenticate() models.Tokens {

	authConfig := auth.DefaultConfig
	// client id from your setup
	authConfig.ClientID = "92aa3738-61c4-44e6-b8ed-4e7a0936e8fc"
	//tenant 2eb75ec2-bed1-41fd-8016-a4e0a0a7072f
	// client secret from your setup
	authConfig.ClientSecret = "PUT_CLIENT_SECRET_HERE"

	// Preform one time login
	authCode := auth.LoginRequest(authConfig)

	t, err := auth.GetTokens(authConfig, authCode, "https://management.azure.com/.default")

	if err != nil {
		fmt.Sprintf("Error")
		panic(err)
	}
	// fmt.Println(t.AccessToken)

	// t contains refresh/access tokens
	fmt.Sprintf("Token: %s", t.AccessToken)
	// getAllSubscriptions(t.AccessToken)
	// getAllTenants(t.AccessToken)
	return t

}

func GetAllTenantsMap(accessToken string) map[string]string {
	url := "https://management.azure.com/tenants?api-version=2020-01-01"

	var m = make(map[string]string)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + accessToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var responseObject models.TenantListResponse
	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.TenantListResponseValue); i++ {
		m[responseObject.TenantListResponseValue[i].TenantName] = responseObject.TenantListResponseValue[i].TenantId
	}

	return m

}

func GetAllSubscriptions(accessToken string) map[string]string {
	url := "https://management.azure.com/subscriptions?api-version=2020-01-01"

	var m = make(map[string]string)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + accessToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var responseObject models.SubscriptionListResponse
	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.SubscriptionListResponseValue); i++ {
		// fmt.Println(responseObject.SubscriptionListResponseValue[i].SubscriptionName)
		m[responseObject.SubscriptionListResponseValue[i].SubscriptionName] = responseObject.SubscriptionListResponseValue[i].SubscriptionId
		// getResourceGroupsPerSubscription(accessToken, responseObject.SubscriptionListResponseValue[i].SubscriptionId)
	}

	return m
}

func GetResourceGroupsPerSubscription(accessToken string, subscriptionId string) map[string]string {
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%v/resourcegroups?api-version=2021-04-01", subscriptionId)

	var m = make(map[string]string)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + accessToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var responseObject models.ResourceGroupResponse
	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.ResourceGroupListResponseValue); i++ {
		m[responseObject.ResourceGroupListResponseValue[i].ResourceGroupName] = responseObject.ResourceGroupListResponseValue[i].Id
	}

	return m

}

func CreateResourceGroup(accessToken string, subscriptionId string, rgName string) {
	url1 := fmt.Sprintf(
		"https://management.azure.com/subscriptions/%v/resourcegroups/%v?api-version=2021-04-01", subscriptionId, rgName)

	fmt.Println(url1)

	// get request URL
	reqURL, _ := url.Parse(url1)

	var bearer = "Bearer " + accessToken

	// create request body
	reqBody := ioutil.NopCloser(strings.NewReader(`
		{
			"location": "eastus"
		}
	`))

	// create a request object
	req := &http.Request{
		Method: "PUT",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json; charset=UTF-8"},
			"Authorization": {bearer},
		},
		Body: reqBody,
	}

	// send an HTTP request using `req` object
	res, err := http.DefaultClient.Do(req)

	// check for response error
	if err != nil {
		log.Fatal("Error:", err)
	}
	// close response body
	res.Body.Close()

	// print response status and body
	if res.StatusCode == 201 {
		fmt.Printf("Resource Group %s Created Succesfully", rgName)
	} else {
		fmt.Printf("Error while creating Resource Group %s", rgName)
	}

}

func CreateAzureApplication(accessToken string, repoName string) string {

	url := "https://graph.microsoft.com/v1.0/applications"
	values := map[string]string{"DisplayName": repoName}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	var bearer = "Bearer " + accessToken
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// print response status and body
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", resp.Body)

	if resp.StatusCode == 201 {
		fmt.Printf("Azure Application with name %s Created Succesfully", repoName)
	} else {
		fmt.Printf("Error while creating Azure Application %s", repoName)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseData))
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var responseObject models.AzureApplicationResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.AppId

}

func CreateServicePrincipal(accessToken string, appId string) {

	url := "https://graph.microsoft.com/v1.0/servicePrincipals"
	values := map[string]string{"appId": appId}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	var bearer = "Bearer " + accessToken
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// print response status and body
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 201 {
		fmt.Printf("SP with AppID %s Created Succesfully", appId)
	} else {
		fmt.Printf("Error while creating SP %s", appId)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseData))

}

func CallOBO(accessToken string) {

	url := "https://login.microsoftonline.com/456742a9-7412-4fb6-881c-3691b12de519/oauth2/v2.0/token"
	values := map[string]string{"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer", "client_id": "92aa3738-61c4-44e6-b8ed-4e7a0936e8fc", "client_secret": "PUT_CLIENT_SECRET_HERE",
		"assertion": accessToken, "scope": "Application.ReadWrite.All", "requested_token_use": "on_behalf_of"}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	// var bearer = "Bearer " + accessToken
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// print response status and body
	fmt.Println(resp.StatusCode)
	// if resp.StatusCode == 201 {
	// 	fmt.Printf("SP with AppID %s Created Succesfully", appId)
	// } else {
	// 	fmt.Printf("Error while creating SP %s", appId)
	// }

	responseData, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseData))

}

type TenantRequest struct {
	ClientId     string `json:"client_id"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	ClientSecret string `json:"client_secret"`
	Grant_Type   string `json:"grant_type"`
}

// GetTokens retrieves access and refresh tokens for a given scope
func GetTokenWithTenantScope(refreshToken string) models.Tokens {
	fmt.Println(refreshToken)
	formVals := url.Values{}
	formVals.Set("grant_type", "refresh_token")
	formVals.Set("refresh_token", refreshToken)
	formVals.Set("scope", "https://graph.microsoft.com/Application.ReadWrite.All")
	formVals.Set("client_secret", "PUT_CLIENT_SECRET_HERE")
	formVals.Set("client_id", "92aa3738-61c4-44e6-b8ed-4e7a0936e8fc")

	response, err := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", formVals)

	// if err != nil {
	// 	return t, errors.Wrap(err, "error while trying to get tokens")
	// }
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println("Tokenssss " + string(body))

	if err != nil {
		// 	return t, errors.Wrap(err, "error while trying to read token json body")
	}

	tokens := models.Tokens{}

	err = json.Unmarshal(body, tokens)
	if err != nil {
		// return t, errors.Wrap(err, "error while trying to parse token json body")
	}

	return tokens

	// fmt.Println(body)

	// tr := TenantRequest{
	// 	ClientId:     "92aa3738-61c4-44e6-b8ed-4e7a0936e8fc",
	// 	RefreshToken: refreshToken,
	// 	Grant_Type:   "refresh_token",
	// 	ClientSecret: "PUT_CLIENT_SECRET_HERE",
	// 	Scope:        "Application.ReadWrite.All",
	// }
	// body, _ := json.Marshal(tr)

	// fmt.Println(string(body))

	// resp, err := http.Post("https://login.microsoftonline.com/2eb75ec2-bed1-41fd-8016-a4e0a0a7072f/oauth2/v2.0/token", "application/x-www-form-urlencoded", bytes.NewBuffer(body))

	// // An error is returned if something goes wrong
	// if err != nil {
	// 	panic(err)
	// }
	// //Need to close the response stream, once response is read.
	// //Hence defer close. It will automatically take care of it.
	// defer resp.Body.Close()

	// body, err1 := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	//Failed to read response.
	// 	panic(err1)
	// }

	// //Convert bytes to String and print
	// jsonStr := string(body)
	// fmt.Println("Response: ", jsonStr)

	// fmt.Println(resp.StatusCode)

	// // //Check response code, if New user is created then read response.
	// // if resp.StatusCode == http.StatusCreated {
	// // 	body, err := ioutil.ReadAll(resp.Body)
	// // 	if err != nil {
	// // 		//Failed to read response.
	// // 		panic(err)
	// // 	}

	// // 	//Convert bytes to String and print
	// // 	jsonStr := string(body)
	// // 	fmt.Println("Response: ", jsonStr)

	// // } else {
	// // 	//The status is not Created. print the error.
	// // 	fmt.Println("Get failed with error: ", resp.Status)
	// // }

}
