package order

import (
	"context"
	"egw-be/sofija/core/ports"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type EgwOrderHttpHandler struct {
	orderSvc ports.EgwOrderUsecase
}

func NewEgwOrderHandler(orderSvc ports.EgwOrderUsecase, wsCont *restful.Container) *EgwOrderHttpHandler {
	httpHandler := &EgwOrderHttpHandler{
		orderSvc: orderSvc,
	}

	ws := new(restful.WebService)

	ws.Path("/order").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/insert").To(httpHandler.InsertOrder))
	ws.Route(ws.DELETE("/delete/{id}").To(httpHandler.DeleteOrder))
	ws.Route(ws.PUT("/update/{id}").To(httpHandler.UpdateOrder))

	wsCont.Add(ws)

	return httpHandler
}

// Performs insert order
func (e *EgwOrderHttpHandler) InsertOrder(req *restful.Request, resp *restful.Response) {
	var reqData InsertRequestData
	req.ReadEntity(&reqData)

	var egwOrder *EgwOrderModel = &EgwOrderModel{}
	egwOrder.UserID = reqData.UserID
	egwOrder.Status = reqData.Status
	egwOrder.CreatedAt = time.Now()
	egwOrder.UpdatedAt = time.Now()

	// Convert the OrderItems from the request to EgwItemOrderModel and add them to the order
	egwOrder.Items = make([]*EgwItemOrderModel, len(reqData.Items))
	for i, item := range reqData.Items {
		egwOrder.Items[i] = &EgwItemOrderModel{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
	}

	// Insert the order into the database
	err := e.insertOrder(req.Request.Context(), egwOrder)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error inserting order"))
		return
	}

	err = generatePDF(egwOrder)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error generating PDF"))
		return
	}

	// Send the inserted order back in the response
	respData := InsertResponseData{
		ID:        egwOrder.ID,
		UserID:    egwOrder.UserID,
		Status:    egwOrder.Status,
		CreatedAt: egwOrder.CreatedAt,
		UpdatedAt: egwOrder.UpdatedAt,
		Items:     make([]*InsertOrderItemResponse, len(egwOrder.Items)),
	}

	// Convert EgwItemOrderModel to OrderItemResponse
	for i, item := range egwOrder.Items {
		respData.Items[i] = &InsertOrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
	}

	resp.WriteAsJson(respData)
}

func (e *EgwOrderHttpHandler) insertOrder(ctx context.Context, egwOrder *EgwOrderModel) error {
	// Convert EgwOrderModel to domain.EgwOrder
	domainOrder := egwOrder.ToDomain()

	// Insert the order into the database using the order service
	insertedOrderID, err := e.orderSvc.InsertOrder(ctx, domainOrder)
	if err != nil {
		return err
	}

	// After successful insertion, retrieve the inserted order from the database to populate it
	fmt.Println(domainOrder.ID)
	insertedOrder, err := e.orderSvc.FindByID(ctx, insertedOrderID)
	if err != nil {
		return err
	}

	// Convert domain.EgwOrder back to EgwOrderModel to update the provided egwOrder with the database data
	egwOrder.FromDomain(insertedOrder)

	return nil
}

func generatePDF(order *EgwOrderModel) error {

	//create new pdf
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	//set font Arial, bold, 16mm for header
	pdf.SetFont("Times", "B", 16)

	// Ispis ID porudžbine i podataka o korisniku
	pdf.Cell(0, 10, fmt.Sprintf("Order ID: %s", order.ID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("User ID: %s", order.UserID))
	pdf.Ln(8)

	//set font and size for table
	pdf.SetFont("Times", "", 12)

	//func drawTable for order items
	drawTable(pdf, order.Items)

	//save pdf
	fileName := fmt.Sprintf("pdf/order/order_%s.pdf", order.ID)
	fmt.Println(fileName)
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	fmt.Printf("PDF %s je uspesno generisan.\n", fileName)
	return nil
}

func drawTable(pdf *gofpdf.Fpdf, items []*EgwItemOrderModel) {

	pdf.SetFillColor(240, 240, 240)

	columnWidths := []float64{10, 50, 20, 25, 35}

	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(columnWidths[0], 10, "#", "1", 0, "C", true, 0, "")
	pdf.CellFormat(columnWidths[1], 10, "Product Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(columnWidths[2], 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(columnWidths[3], 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(columnWidths[4], 10, "Total Price", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)

	for i, item := range items {
		fmt.Println(item.ProductName)
		fmt.Println(item.Quantity)
		fmt.Println(item.Price)

		pdf.CellFormat(columnWidths[0], 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[1], 10, item.ProductName, "1", 0, "", false, 0, "")
		pdf.CellFormat(columnWidths[2], 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(columnWidths[3], 10, fmt.Sprintf("%.2f", item.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(columnWidths[4], 10, fmt.Sprintf("%.2f", float64(item.Quantity)*item.Price), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
}

func (e *EgwOrderHttpHandler) DeleteOrder(req *restful.Request, resp *restful.Response) {
	// get order ID for delete query
	reqID := req.PathParameter("id")
	fmt.Println(reqID)
	if len(reqID) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no ID provided for order deletion"))
		return
	}

	_, err := uuid.Parse(reqID)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid ID format"))
		return
	}

	ctx := req.Request.Context()

	err = e.orderSvc.Delete(ctx, reqID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error deleting order"))
		return
	}

	resp.WriteHeader(http.StatusOK)
}

func (e *EgwOrderHttpHandler) UpdateOrder(req *restful.Request, resp *restful.Response) {

	var a UpdateRequestData
	req.ReadEntity(&a)

	orderID := req.PathParameter("id")
	if len(orderID) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no id found for order"))
		return
	}

	ctx := req.Request.Context()
	dataOrder := &domain.EgwOrder{ID: orderID, Status: a.Status}

	err := e.orderSvc.Update(ctx, dataOrder)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error updating order"))
		return
	}

	returnedOrder, err := e.orderSvc.FindByID(ctx, orderID)
	if err != nil {
		return
	}

	// return updated order as data
	resp.WriteAsJson(returnedOrder)
}
