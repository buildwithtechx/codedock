package support

type CloudSupportCategory string

const (
	CloudSupportCategoryDeployment CloudSupportCategory = "deployment"
	CloudSupportCategoryBilling    CloudSupportCategory = "billing"
	CloudSupportCategoryAccount    CloudSupportCategory = "account"
	CloudSupportCategoryGeneral    CloudSupportCategory = "general"
)

type CloudSupportStatus string

const (
	CloudSupportStatusOpen     CloudSupportStatus = "open"
	CloudSupportStatusResolved CloudSupportStatus = "resolved"
)

type CloudSupportAuthor string

const (
	CloudSupportAuthorCustomer CloudSupportAuthor = "customer"
	CloudSupportAuthorSupport  CloudSupportAuthor = "support"
)
