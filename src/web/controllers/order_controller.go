package controllers

import (
	"context"
	"demoService/src/application/services"
	"demoService/src/domain"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type OrderController struct {
	orderService services.OrderService
}

func NewOrderController(orderService services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (controller *OrderController) UseHandlers() {
	http.HandleFunc("/order/", controller.GetOrderById)
	http.HandleFunc("/", controller.OrderView)
}

// GetOrderById godoc
// @Summary      Получить заказ по ID
// @Description  Возвращает заказ по его уникальному идентификатору
// @Tags         orders
// @Param        order_uid   path      string  true  "Уникальный ID заказа"
// @Produce      json
// @Success      200  {object}  models.Order
// @Failure      400  {string}  string  "order_uid not provided"
// @Failure      404  {string}  string  "Order not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /order/{order_uid} [get]
func (controller *OrderController) GetOrderById(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(response, "order_uid not provided", http.StatusBadRequest)
		return
	}

	orderId := parts[2]
	order, err := controller.orderService.GetOrderById(request.Context(), orderId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.NotFound(response, request)
			return
		case errors.Is(err, context.Canceled):
			log.Println("request canceled by client")
			return
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("request deadline exceeded")
			return
		default:
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	_, err = response.Write(jsonBytes)
	if err != nil {
		log.Println(err.Error())
	}
}

func (controller *OrderController) OrderView(response http.ResponseWriter, _ *http.Request) {
	view, err := template.ParseFiles("src/web/templates/order.html")
	if err != nil {
		http.Error(response, "Template not found: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.Execute(response, nil); err != nil {
		http.Error(response, "Render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
