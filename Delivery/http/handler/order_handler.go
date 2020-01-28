package handler

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
)

// OrderHandler handles order handler admin requests
type OrderHandler struct {
	tmpl        *template.Template
	orderServ   order.OrderService
	usdrServ    user.UserService
	itemServ    productpage.ItemService
	csrfSignKey []byte
}

var idr string
// NewOrderHandler initializes and returns new OrderHandler
func NewOrderHandler(t *template.Template, os order.OrderService, us user.UserService,
	is productpage.ItemService, csKey []byte) *OrderHandler {
	return &OrderHandler{tmpl: t, orderServ: os, usdrServ: us, itemServ: is, csrfSignKey: csKey}
}

// Orders handle requests on route /orders
func (oh *OrderHandler) Orders(w http.ResponseWriter, r *http.Request) {
	orders, _ := oh.orderServ.Orders()
	token, err := csrfToken.CSRFToken(oh.csrfSignKey)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Orders  []entity.Order
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Orders:  orders,
		CSRF:    token,
	}
	err = oh.tmpl.ExecuteTemplate(w, "admin.order.layout", tmplData)
	if err != nil {
		panic(err.Error())
	}
}


func (oh *OrderHandler) GetSingleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		ordr, errs := oh.orderServ.Order(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		err = oh.tmpl.ExecuteTemplate(w, "admin.order.layout", ordr)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	http.Redirect(w, r, "/admin/order", 303)
}
func (oh *OrderHandler) GetUserOrder(w http.ResponseWriter, r *http.Request) {
	user := &entity.User{}
	idRaw := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idRaw)
	products := []entity.Product{}
	user.ID = uint(id)
	// var id uint
	order, _ := oh.orderServ.CustomerOrders(user)
	productlist := strings.Split(order.ItemsID, ",")
	for i:=0; i< len(productlist); i++{
		productid := productlist[i]
		prodid, _ := strconv.Atoi(productid)
		pro, _ := oh.itemServ.Item(uint(prodid))
		products = append(products, *pro)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Order   entity.Order
		Products []entity.Product
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Order:   order,
		Products: products,
		// CSRF:    token,
	}
	err := oh.tmpl.ExecuteTemplate(w, "checkoutpage.html", tmplData)
	if err != nil {
		panic(err.Error())
	}
}

func (oh *OrderHandler) OrderNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid, _ := strconv.Atoi(r.FormValue("userid"))
		t := time.Now()
		total, _ := strconv.ParseFloat(r.FormValue("total"), 64)
		//prodids, _ := strconv.Atoi(r.FormValue("prodids"))
		prodids := r.FormValue("prodids")
		log.Println("User order id: ", uid)
		log.Println("User order prod ids: ", prodids)
		log.Println("User order total: ", total)
		idr = strconv.Itoa(uid)
		ord := &entity.Order{
			UserID:    uint(uid),
			CreatedAt: t,
			ItemsID:   prodids,
			Total:     total,
		}
		// link += string(pid)
		_, errs := oh.orderServ.StoreOrder(ord)
		if len(errs) > 0 {
			panic(errs)
		}
	}
	http.Redirect(w, r, "/getorder?id="+idr, 303)
}

// OrderDelete handles requests on route /order/delete
func (oh *OrderHandler) OrderDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := oh.orderServ.DeleteOrder(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		http.Redirect(w, r, "/getusercart?id="+idr, http.StatusSeeOther)
	}
}
