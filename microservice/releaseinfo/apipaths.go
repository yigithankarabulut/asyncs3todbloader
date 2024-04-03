package releaseinfo

const (
	apiVersion = "/api/" + Version
	Product    = apiVersion + "/product"
)

const (
	GetProduct = Product + "/:id"
)
