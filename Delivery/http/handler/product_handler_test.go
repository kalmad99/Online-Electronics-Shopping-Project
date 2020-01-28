package handler

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/repository"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/service"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
)

func TestSellerProduct(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/seller/products", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	url := ts.URL

	resp, err := tc.Get(url + "/seller/products")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestproductCategoriesNew(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/seller/products/new", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}
	form.Add("name", entity.ProductMock.Name)
	form.Add("Description", entity.ProductMock.Description)
	form.Add("Image", entity.ProductMock.Image)
	form.Add("CategoryID", string(entity.ProductMock.CategoryID))
	form.Add("Quantity", string(entity.ProductMock.Quantity))
	//form.Add("Price",string(entity.ProductMock.Price))
	//form.Add("rate",entity.ProductMock.Rating))

	resp, err := tc.PostForm(sURL+"/seller/products/new", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestSellerProductsUpdate(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/seller/products/update", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.ProductMock.CategoryID))
	form.Add("Name", entity.ProductMock.Name)
	form.Add("description", entity.ProductMock.Description)
	form.Add("Image", entity.ProductMock.Image)

	resp, err := tc.PostForm(sURL+"/seller/products/update?id=1", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestSellerProductsDelete(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/seller/products/delete", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.ProductMock.CategoryID))
	form.Add("Name", entity.ProductMock.Name)
	form.Add("description", entity.ProductMock.Description)
	form.Add("Image", entity.ProductMock.Image)

	resp, err := tc.PostForm(sURL+"/seller/products/delete?id=1", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}

func TestSearchProducts(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/searchProducts", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.ProductMock.CategoryID))
	form.Add("Name", entity.ProductMock.Name)
	form.Add("description", entity.ProductMock.Description)
	form.Add("Image", entity.ProductMock.Image)

	resp, err := tc.PostForm(sURL+"/searchProducts?id=1", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}
func TestProductDetail(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockProductRepo(nil)
	categoryServv := service.NewItemService(categoryRepo)

	sellerCtgHandler := NewSellerProductHandler(tmpl, categoryServv, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/detail", sellerCtgHandler.SellerProducts)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.ProductMock.CategoryID))
	form.Add("Name", entity.ProductMock.Name)
	form.Add("description", entity.ProductMock.Description)
	form.Add("Image", entity.ProductMock.Image)

	resp, err := tc.PostForm(sURL+"/detail?id=1", form)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(body, []byte("name")) {
		t.Errorf("want body to contain %q", body)
	}

}
