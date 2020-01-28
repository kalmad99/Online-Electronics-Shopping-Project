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
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
)

// AdminCategoryHandler handles category handler admin requests
type CartHandler struct {
	tmpl           *template.Template
	cartSrv        cart.CartService
	userService    user.UserService
	sessionService user.SessionService
	userSess       *entity.Session
	loggedInUser   *entity.User
	userRole       user.RoleService
	productService productpage.ItemService
	csrfSignKey    []byte
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewCartHandler(t *template.Template, cs cart.CartService, usrServ user.UserService,
	sessServ user.SessionService, uRole user.RoleService, usrSess *entity.Session, prodServ productpage.ItemService,
	csKey []byte) *CartHandler {
	return &CartHandler{tmpl: t, cartSrv: cs, userService: usrServ, sessionService: sessServ,
		userRole: uRole, userSess: usrSess, productService: prodServ, csrfSignKey: csKey}
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

func (ch *CartHandler) AddtoCart(w http.ResponseWriter, r *http.Request) {

}

func (ch *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {

}