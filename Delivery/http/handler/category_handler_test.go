package handler

import (
	"bytes"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/repository"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/service"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAdminCategories(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockCategoryRepo(nil)
	categoryServ := service.NewCategoryService(categoryRepo)

	adminCatgHandler := NewAdminCategoryHandler(tmpl, categoryServ, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories", adminCatgHandler.AdminCategories)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	url := ts.URL

	resp, err := tc.Get(url + "/admin/categories")
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

func TestAdminCategoriesNew(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockCategoryRepo(nil)
	categoryServ := service.NewCategoryService(categoryRepo)

	adminCatgHandler := NewAdminCategoryHandler(tmpl, categoryServ, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories/new", adminCatgHandler.AdminCategories)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}
	form.Add("name", entity.CategoryMock.Name)
	form.Add("Description", entity.CategoryMock.Description)
	form.Add("Image", entity.CategoryMock.Image)

	resp, err := tc.PostForm(sURL+"/admin/categories/new", form)
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

func TestAdminCategoresUpdate(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockCategoryRepo(nil)
	categoryServ := service.NewCategoryService(categoryRepo)

	adminCatgHandler := NewAdminCategoryHandler(tmpl, categoryServ, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories/update", adminCatgHandler.AdminCategories)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.CategoryMock.ID))
	form.Add("Name", entity.CategoryMock.Name)
	form.Add("description", entity.CategoryMock.Description)
	form.Add("Image", entity.CategoryMock.Image)

	resp, err := tc.PostForm(sURL+"/admin/categories/update?id=1", form)
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

func TestAdminCategoresDelete(t *testing.T) {

	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))

	categoryRepo := repository.NewMockCategoryRepo(nil)
	categoryServ := service.NewCategoryService(categoryRepo)

	adminCatgHandler := NewAdminCategoryHandler(tmpl, categoryServ, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories/delete", adminCatgHandler.AdminCategories)
	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	tc := ts.Client()
	sURL := ts.URL

	form := url.Values{}

	form.Add("ID", string(entity.CategoryMock.ID))
	form.Add("Name", entity.CategoryMock.Name)
	form.Add("description", entity.CategoryMock.Description)
	form.Add("Image", entity.CategoryMock.Image)

	resp, err := tc.PostForm(sURL+"/admin/categories/delete?id=1", form)
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
