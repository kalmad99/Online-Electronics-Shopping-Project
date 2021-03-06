package handler

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var outid, categoid uint

// ProductHandler handles product handler admin requests
type SellerProductHandler struct {
	tmpl        *template.Template
	productSrv  productpage.ItemService
	csrfSignKey []byte
}

// NewSellerProductHandler initializes and returns new SellerProductHandler
func NewSellerProductHandler(t *template.Template, is productpage.ItemService, csKey []byte) *SellerProductHandler {
	return &SellerProductHandler{tmpl: t, productSrv: is, csrfSignKey: csKey}
}

// SellerProducts handle requests on route /admin/products
func (sph *SellerProductHandler) SellerProducts(w http.ResponseWriter, r *http.Request) {
	products, errs := sph.productSrv.Items()
	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
	if err != nil {
		panic(errs)
	}
	tmplData := struct {
		Values   url.Values
		VErrors  form.ValidationErrors
		Products []entity.Product
		CSRF     string
	}{
		Values:   nil,
		VErrors:  nil,
		Products: products,
		CSRF:     token,
	}

	err = sph.tmpl.ExecuteTemplate(w, "seller.products.layout", tmplData)
	if err != nil {
		panic(err.Error())
	}
}

func (sph *SellerProductHandler) SellerProductsNew(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		newProForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		err := sph.tmpl.ExecuteTemplate(w, "seller.product.new.layout", newProForm)
		if err != nil {
			panic(err.Error())
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		newProForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newProForm.MinLength("description", 10)
		newProForm.CSRF = token

		mf, fh, err := r.FormFile("catimg")
		if err != nil {
			newProForm.VErrors.Add("catimg", "File error")
			err := sph.tmpl.ExecuteTemplate(w, "seller.product.new.layout", newProForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		defer mf.Close()

		pro := &entity.Product{}
		pro.Name = r.FormValue("name")
		pro.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		pro.Description = r.FormValue("description")
		pro.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
		categ, _ := strconv.Atoi(r.FormValue("type"))
		pro.CategoryID = uint(categ)
		pro.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		_, errs := sph.productSrv.StoreItem(pro)

		if len(errs) > 0 {
			panic(errs)
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		pro1 := &entity.Product{CategoryID: pro.CategoryID, ID: pro.ID}
		errs = sph.productSrv.StoreItemCateg(pro1)
		if len(errs) > 0 {
			panic(errs)
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, "/admin/products", http.StatusSeeOther)

		}
	}
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (sph *SellerProductHandler) SellerProductsUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}
		//if err != nil {
		//	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		//}

		pro, errs := sph.productSrv.Item(uint(id))

		outid = uint(id)
		log.Println("outid", outid)
		if len(errs) > 0 {
			panic(errs)
		}
		//if len(errs) > 0 {
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//}
		price := strconv.FormatFloat(pro.Price, 'f', 2, 64)
		rating := strconv.FormatFloat(pro.Rating, 'f', 2, 64)
		ratercount := strconv.FormatFloat(pro.RatersCount, 'f', 2, 64)
		quan := strconv.Itoa(pro.Quantity)
		catid := pro.CategoryID
		categoid = catid
		values := url.Values{}
		values.Add("proid", idRaw)
		values.Add("name", pro.Name)
		//values.Add("type", cate)
		values.Add("description", pro.Description)
		values.Add("price", price)
		values.Add("quantity", quan)
		values.Add("catimg", pro.Image)
		values.Add("ratcount", ratercount)
		values.Add("rate", rating)
		upProForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			Product *entity.Product
			CSRF    string
		}{
			Values:  values,
			VErrors: form.ValidationErrors{},
			Product: pro,
			CSRF:    token,
		}

		err = sph.tmpl.ExecuteTemplate(w, "seller.products.update.layout", upProForm)
		if err != nil {
			err.Error()
		}
		return
	}
	if r.Method == http.MethodPost {

		log.Println("ID", outid)
		if err != nil {
			panic(err.Error())
		}
		quan, _ := strconv.Atoi(r.FormValue("quantity"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		rating, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
		ratercount, _ := strconv.ParseFloat(r.FormValue("ratcount"), 64)
		prod := &entity.Product{
			ID:          outid,
			Name:        r.FormValue("name"),
			CategoryID:  categoid,
			Description: r.FormValue("description"),
			Quantity:    quan,
			Price:       price,
			RatersCount: ratercount,
			Rating:      rating,
			Image:       r.FormValue("imgname"),
		}

		log.Println("Name", prod.Name)
		log.Println("Price", prod.Price)
		log.Println("Descr", prod.Description)
		log.Println("Quan", prod.Quantity)
		log.Println("Image", prod.Image)
		log.Println("rate", prod.Rating)
		log.Println("count", prod.RatersCount)

		mf, fh, err := r.FormFile("catimg")
		if err == nil {
			prod.Image = fh.Filename
			err = writeFile(&mf, prod.Image)
		}
		if mf != nil {
			defer mf.Close()
		}

		_, errs := sph.productSrv.UpdateItem(prod)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}
}

// AdminCategoriesDelete handle requests on route /admin/categories/delete
func (sph *SellerProductHandler) SellerProductsDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		_, errs := sph.productSrv.DeleteItem(uint(id))

		if len(errs) > 0 {
			panic(err)
		}

	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (sph *SellerProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		res := r.URL.Query().Get("search")

		if len(res) == 0 {
			http.Redirect(w, r, "/", 303)
		}
		results, err := sph.productSrv.SearchProduct(res)

		if err != nil {
			panic(err)
		}

		sph.tmpl.ExecuteTemplate(w, "searchresults.layout", results)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (sph *SellerProductHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		pro, errs := sph.productSrv.Item(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		_ = sph.tmpl.ExecuteTemplate(w, "productdetail.layout", pro)
	}
}

func (sph *SellerProductHandler) Rating(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		idRaw := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idRaw)

		pro, errs := sph.productSrv.Item(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}
		_ = sph.tmpl.ExecuteTemplate(w, "ratings.html", pro)
	} else if req.Method == http.MethodPost {

		prod := &entity.Product{}
		idRaw, _ := strconv.Atoi(req.FormValue("id"))
		prod.ID = uint(idRaw)
		prod.Rating, _ = strconv.ParseFloat(req.FormValue("star"), 64)

		log.Println("prod.rating", prod.Rating)
		log.Println("prod.id", prod.ID)
		_, err := sph.productSrv.RateProduct(prod)

		if err != nil {
			panic(err)
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)

	} else {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
