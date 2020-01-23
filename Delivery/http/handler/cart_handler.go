package handler

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// AdminCategoryHandler handles category handler admin requests
type CartHandler struct {
	tmpl        *template.Template
	cartSrv     cart.CartService
	csrfSignKey []byte
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewCartHandler(t *template.Template, cs cart.CartService, csKey []byte) *CartHandler {
	return &CartHandler{tmpl: t, cartSrv: cs, csrfSignKey: csKey}
}

// Products handle requests on route /seller/products
func (ch *CartHandler) GetCarts(w http.ResponseWriter, r *http.Request) {
	carts, _ := ch.cartSrv.GetCarts()
	token, err := csrfToken.CSRFToken(ch.csrfSignKey)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Carts   []entity.Cart
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Carts:   carts,
		CSRF:    token,
	}
	err = ch.tmpl.ExecuteTemplate(w, "admin.cart.detail.layout", tmplData)
	if err != nil {
		panic(err.Error())
	}
}

func (ch *CartHandler) GetSingleCart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		cart, errs := ch.cartSrv.GetSingleCart(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		err = ch.tmpl.ExecuteTemplate(w, "seller.products.update.layout", cart)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	http.Redirect(w, r, "/admin/cart", 303)
}

func (ch *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {

}

func (ch *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {

}

func (ch *CartHandler) AddtoCart(w http.ResponseWriter, r *http.Request) {

}

func (ch *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {

}

///////////////////////////////////////////////////////////////////////////////////////////////////
// AdminCategories handle requests on route /admin/categories
//func (ch *CartHandler) ItemsinCart(w http.ResponseWriter, r *http.Request) {
//	token, err := csrfToken.CSRFToken(ch.csrfSignKey)
//	if err != nil {
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//	}
//	if r.Method == http.MethodGet {
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//		if err != nil {
//			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		cart, errs := ch.cartSrv.GetCart(uint(id))
//		if len(errs) > 0 {
//			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//			return
//		}
//		values:=url.Values{}
//		values.Add("cartid", idRaw)
//		values.Add("userid", string(cart.UserID))
//		values.Add("itemid", string(cart.Items[0].ID))
//		values.Add("quantity", string(cart.Items[0].Quantity))
//		values.Add("placedat", cart.PlacedAt.Format("2006-01-02 15:04:05"))
//
//		itmCart := struct {
//			Values  url.Values
//			VErrors form.ValidationErrors
//			Cart    *entity.Cart
//			CSRF    string
//		}{
//			Values:  values,
//			VErrors: form.ValidationErrors{},
//			Cart:    cart,
//			CSRF:    token,
//		}
//		ch.tmpl.ExecuteTemplate(w, "admin.user.update.layout", itmCart)
//		return
//	}
//}
//
//func (ch *CartHandler) AddtoCart(w http.ResponseWriter, r *http.Request) {
//	token, err := csrfToken.CSRFToken(ch.csrfSignKey)
//	if err != nil {
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//	}
//	if r.Method == http.MethodGet {
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//		if err != nil {
//			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//
//		if len(errs) > 0 {
//			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//			return
//		}
//		values := url.Values{}
//		values.Add("cartid", idRaw)
//		values.Add("userid", string(cart.UserID))
//		values.Add("itemid", string(cart.ItemID))
//		values.Add("quantity", string(cart.Quantity))
//		values.Add("placedat", cart.PlacedAt.Format("2006-01-02 15:04:05"))
//
//		itmCart := struct {
//			Values  url.Values
//			VErrors form.ValidationErrors
//			Cart    *entity.Cart
//			CSRF    string
//		}{
//			Values:  values,
//			VErrors: form.ValidationErrors{},
//			Cart:    cart,
//			CSRF:    token,
//		}
//		ch.tmpl.ExecuteTemplate(w, "admin.user.update.layout", itmCart)
//		return
//	}
//}
