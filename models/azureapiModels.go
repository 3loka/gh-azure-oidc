package models

type SubscriptionListResponse struct {
	SubscriptionListResponseValue []SubscriptionListResponseValue `json:"value"`
}

type SubscriptionListResponseValue struct {
	Id               string `json:"id"`
	SubscriptionId   string `json:"subscriptionId"`
	SubscriptionName string `json:"displayName"`
}

type TenantListResponse struct {
	TenantListResponseValue []TenantListResponseValue `json:"value"`
}

type TenantListResponseValue struct {
	Id         string `json:"id"`
	TenantId   string `json:"tenantId"`
	TenantName string `json:"displayName"`
}

type ResourceGroupResponse struct {
	ResourceGroupListResponseValue []ResourceGroupListResponseValue `json:"value"`
}

type ResourceGroupListResponseValue struct {
	Id                string `json:"id"`
	ResourceGroupName string `json:"name"`
}

type AzureApplicationResponse struct {
	Id    string `json:"id"`
	AppId string `json:"appId"`
}

type ServicePrincipalResponse struct {
	Id                    string `json:"id"`
	AppId                 string `json:"appId"`
	ServicePrincipalNames string `json:"servicePrincipalNames"`
	DisplayName           string `json:"displayName"`
}

type CreateResourceGroupResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
