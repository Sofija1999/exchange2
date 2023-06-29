package product

import (
	"context"
	"fmt"
	"net/http"

	//"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"errors"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/emicklei/go-restful/v3"
	//"log"
)

type EgwProductHttpHandler struct {
	productSvc ports.EgwProductUsecase
}

func NewEgwProductHandler(productSvc ports.EgwProductUsecase, wsCont *restful.Container) *EgwProductHttpHandler {
	httpHandler := &EgwProductHttpHandler{
		productSvc: productSvc,
	}

	ws := new(restful.WebService)

	ws.Path("/product").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/insert").To(httpHandler.InsertProduct))
	ws.Route(ws.PUT("/update/id").To(httpHandler.UpdateProduct))
	ws.Route(ws.DELETE("/delete/id").To(httpHandler.DeleteProduct))
	ws.Route(ws.GET("/get-all").To(httpHandler.GetAllProducts))

	wsCont.Add(ws)

	return httpHandler
}

// Performs insert product
func (e *EgwProductHttpHandler) InsertProduct(req *restful.Request, resp *restful.Response) {
	var reqData InsertRequestData
	req.ReadEntity(&reqData)

	var egwProduct *EgwProductModel = &EgwProductModel{}
	egwProduct.Name = reqData.Name
	egwProduct.ShortDescription = reqData.ShortDescription
	egwProduct.Description = reqData.Description
	egwProduct.Price = reqData.Price

	err := e.insertProduct(req.Request.Context(), egwProduct)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error insert product"))
		return
	}

	// send product back
	respData := InsertResponseData{Product: *egwProduct}

	resp.WriteAsJson(respData)
}

func (e *EgwProductHttpHandler) insertProduct(ctx context.Context, egwProduct *EgwProductModel) error {
	err := e.productSvc.InsertProduct(ctx, egwProduct.ToDomain())
	if err != nil {
		return err
	}
	// retrieve their data from the DB to populate it (e.g. NAME)
	productData, err := e.productSvc.FindByName(ctx, egwProduct.Name)
	if err != nil {
		return err
	}
	egwProduct.FromDomain(productData)
	return nil
}

func (e *EgwProductHttpHandler) UpdateProduct(req *restful.Request, resp *restful.Response) {

	var a UpdateRequestData
	req.ReadEntity(&a)

	productID := req.PathParameter("id")
	if len(productID) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no id found for product"))
		return
	}

	ctx := req.Request.Context()
	dataProduct := &domain.EgwProduct{ID: productID, ShortDescription: a.ShortDescription, Description: a.Description, Price: a.Price}

	err := e.productSvc.Update(ctx, dataProduct)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error updating product"))
		return
	}

	// return updated product as data
	var retProduct *EgwProductModel = &EgwProductModel{}
	retProduct.FromDomain(dataProduct)
	resp.WriteAsJson(retProduct)
}

func (e *EgwProductHttpHandler) DeleteProduct(req *restful.Request, resp *restful.Response) {
	// get product ID for delete query
	reqID := req.PathParameter("id")
	fmt.Println(reqID)
	if len(reqID) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no ID provided for product deletion"))
		return
	}

	ctx := req.Request.Context()

	err := e.productSvc.Delete(ctx, reqID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error deleting product"))
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (e *EgwProductHttpHandler) GetAllProducts(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	products, err := e.productSvc.GetAll(ctx)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error retrieving products"))
		return
	}

	var productModels []EgwProductModel
	for _, product := range products {
		productModel := EgwProductModel{}
		productModel.FromDomain(product) // Pass the pointer to domain.EgwProduct
		productModels = append(productModels, productModel)
	}

	resp.WriteAsJson(productModels)
}
