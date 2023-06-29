package product

type InsertRequestData struct {
	Name             string
	ShortDescription string
	Description      string
	Price            int64
}

type InsertResponseData struct {
	Product EgwProductModel
}

type UpdateRequestData struct {
	ShortDescription string
	Description      string
	Price            int64
}

type UpdateResponseData struct {
	Product EgwProductModel
}
