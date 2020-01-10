package handler

import (
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"../../../entity"
	"../../../productpage"
)

// AdminCategoryHandler handles category handler admin requests
type SellerProductHandler struct {
	tmpl       *template.Template
	productSrv productpage.ProductService
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewSellerProductHandler(t *template.Template, is productpage.ProductService) *SellerProductHandler {
	return &SellerProductHandler{tmpl: t, productSrv: is}
}

// AdminCategories handle requests on route /admin/categories
func (sph *SellerProductHandler) SellerProducts(w http.ResponseWriter, r *http.Request) {
	products, errs := sph.productSrv.Products()
	if errs != nil {
		panic(errs)
	}
	err := sph.tmpl.ExecuteTemplate(w, "seller.products.layout", products)
	if err != nil {
		panic(err.Error())
	}
}

// AdminCategoriesNew hanlde requests on route /admin/categories/new
func (sph *SellerProductHandler) SellerProductsNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		pro := entity.Product{}
		pro.Name = r.FormValue("name")
		pro.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		pro.Description = r.FormValue("description")
		pro.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)

		mf, fh, err := r.FormFile("catimg")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		pro.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		err = sph.productSrv.StoreProduct(pro)

		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)

	} else {

		err := sph.tmpl.ExecuteTemplate(w, "seller.product.new.layout", nil)
		if err != nil {
			panic(err.Error())
		}

	}
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (sph *SellerProductHandler) SellerProductsUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		pro, err := sph.productSrv.Product(id)

		if err != nil {
			panic(err)
		}

		sph.tmpl.ExecuteTemplate(w, "seller.products.update.layout", pro)

	} else if r.Method == http.MethodPost {

		prod := entity.Product{}
		id, _ := strconv.Atoi(r.FormValue("id"))
		prod.ID = id
		prod.Name = r.FormValue("name")
		prod.Description = r.FormValue("description")
		prod.Image = r.FormValue("image")
		mf, _, err := r.FormFile("catimg")

		if err != nil {
			panic(err)
		}

		defer mf.Close()

		writeFile(&mf, prod.Image)

		err = sph.productSrv.UpdateProduct(prod)

		if err != nil {
			panic(err.Error())
		}
		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
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

		err = sph.productSrv.DeleteProduct(id)

		if err != nil {
			panic(err.Error())
		}
	}

	http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
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

		pro, err := sph.productSrv.Product(id)

		if err != nil {
			panic(err.Error)
		}

		_ = sph.tmpl.ExecuteTemplate(w, "productdetail.layout", pro)
	}
}
func (sph *SellerProductHandler) SearchProduct(index string) ([]entity.Product, error) {
	products, err := sph.productSrv.SearchProduct(index)

	if err != nil {
		return nil, err
	}

	return products, nil
}

//func (sph *SellerProductHandler) RateProduct(pro entity.Product) (entity.Product, []error) {
//
//	prowithrate, err := sph.productSrv.RateProduct(pro)
//	if err != nil {
//		return prowithrate, err
//	}
//	return prowithrate, nil
//}

func writeFile(mf *multipart.File, fname string) {

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "../", "../", "ui", "assets", "img", fname)
	image, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	io.Copy(image, *mf)
}
