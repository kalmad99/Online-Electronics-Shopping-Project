package handler

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
)

// AdminCategoryHandler handles category handler admin requests
type CartHandler struct {
	tmpl           *template.Template
	cartSrv        cart.CartService
	userService    user.UserService
	productService productpage.ItemService
	csrfSignKey    []byte
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewCartHandler(t *template.Template, cs cart.CartService, usrServ user.UserService,
	sessServ user.SessionService, uRole user.RoleService, usrSess *entity.Session, prodServ productpage.ItemService,
	csKey []byte) *CartHandler {
	return &CartHandler{tmpl: t, cartSrv: cs, userService: usrServ, productService: prodServ,
		csrfSignKey: csKey}
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
	user := &entity.User{}
	//if r.Method == http.MethodGet{
	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	//if err != nil {
	//	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//}
	user.ID = uint(id)
	prds, _ := ch.cartSrv.GetUserCart(user)
	productsss := []entity.Product{}
	for i := range prds {
		pid := prds[i].ID
		product, _ := ch.productService.Item(pid)
		productsss = append(productsss, *product)
	}

	//token, err := csrfToken.CSRFToken(ch.csrfSignKey)
	//
	//if err != nil {
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//}
	tmplData := struct {
		Values   url.Values
		VErrors  form.ValidationErrors
		Products []entity.Product
		CSRF     string
	}{
		Values:   nil,
		VErrors:  nil,
		Products: productsss,
		//CSRF:    token,
	}
	err = ch.tmpl.ExecuteTemplate(w, "cart.html", tmplData)
	if err != nil {
		panic(err.Error())
		//}
	}
}

func (ch *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		usr, _ := ch.userService.User(uint(id))
		_, errs := ch.cartSrv.DeleteCart(usr)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	//http.Redirect(w, r, "/admin/cart", 303)
}

var idk string
func (ch *CartHandler) AddtoCart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid, _ := strconv.Atoi(r.FormValue("userid"))
		pid, _ := strconv.Atoi(r.FormValue("prodid"))
		t := time.Now()
		pri, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		log.Println("User cart id: ", uid)
		log.Println("User cart prod id: ", pid)
		log.Println("User cart price: ", pri)
		idk = strconv.Itoa(uid)
		car := &entity.Cart{
			UserId:    uint(uid),
			ProductId: uint(pid),
			AddedTime: t,
			Price:     pri,
		}
		// link += string(pid)
		_, errs := ch.cartSrv.AddtoCart(car)
		if len(errs) > 0 {
			panic(errs)
		}
	}
	http.Redirect(w, r, "/getusercart?id="+idk, 303)
}

func (ch *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	cart := entity.Cart{}
	idd, _ := strconv.Atoi(cid)
	cart.UserId = uint(idd)
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		cart.ProductId = uint(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := ch.cartSrv.UpdateCart(&cart)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	http.Redirect(w, r, "/getusercart?id=" + cid, 303)
}
