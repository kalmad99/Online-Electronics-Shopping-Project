package handler

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"html/template"
	"net/http"
	"net/url"
)

// MenuHandler handles menu related requests
type MenuHandler struct {
	tmpl        *template.Template
	productSrv  productpage.ItemRepository
	csrfSignKey []byte
}

// NewMenuHandler initializes and returns new MenuHandler
func NewMenuHandler(T *template.Template, IS productpage.ItemService, csKey []byte) *MenuHandler {
	return &MenuHandler{tmpl: T, productSrv: IS, csrfSignKey: csKey}
}

//Index
func (mh *MenuHandler) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	token, err := csrfToken.CSRFToken(mh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	products, errs := mh.productSrv.Items()
	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmplData := struct {
		Values   url.Values
		VErrors  form.ValidationErrors
		Products []entity.Product
		CSRF     string
		UserID   string
	}{
		Values:   nil,
		VErrors:  nil,
		Products: products,
		CSRF:     token,
		UserID:   r.FormValue("userid"),
	}

	mh.tmpl.ExecuteTemplate(w, "index.layout", tmplData)
}

// Admin handle request on route /admin
func (mh *MenuHandler) Admin(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(mh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		CSRF:    token,
	}
	mh.tmpl.ExecuteTemplate(w, "admin.index.layout", tmplData)
}

func (mh *MenuHandler) RegistPage(w http.ResponseWriter, req *http.Request) {
	mh.tmpl.ExecuteTemplate(w, "Registrationform.layout", nil)
}

func (mh *MenuHandler) LoginPage(w http.ResponseWriter, req *http.Request) {

	mh.tmpl.ExecuteTemplate(w, "login.html", nil)
}

func (mh *MenuHandler) ProductDetail(w http.ResponseWriter, req *http.Request) {
	mh.tmpl.ExecuteTemplate(w, "productdetail.html", nil)
}

func (mh *MenuHandler) PaySuccess(w http.ResponseWriter, req *http.Request) {
	mh.tmpl.ExecuteTemplate(w, "pay.success.html", nil)
}
