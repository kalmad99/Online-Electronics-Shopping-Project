//

package handler

//import (
//	//"Online-Electronics-Shopping-Project3/allEntitiesAction/cart/repository"
//	//"Online-Electronics-Shopping-Project3/allEntitiesAction/productpage/repository"
//	//"Online-Electronics-Shopping-Project3/allEntitiesAction/productpage/service"
//	//csepim	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/service"
//	//crepim	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/repository"
//	//urepim	"Online-Electronics-Shopping-Project3/allEntitiesAction/user/service"
//	//"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/service"
//	//"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/repository"
//
//	//"Online-Electronics-Shopping-Project3/allEntitiesAction/cart/repository"
//	////"Online-Electronics-Shopping-Project3/allEntitiesAction/cart/repository"
//	//prepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/repository"
//	psrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/service"
//
//	urepimp "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user/repository"
//	usrvimp "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user/service"
//
//	crepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/repository"
//	csrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/service"
//
//	//orepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order/repository"
//	//osrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order/usecase"
//	//
//	//brepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank/repository"
//	//bsrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank/usecase"
//
//	"html/template"
//
//	"net/http"
//	"net/http/httptest"
//
//	"testing"
//)
//
//func TestGetCarts(t *testing.T) {
//
//	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))
//
//	categoryRepo := crepim.NewCartGormRepo(nil)
//	categoryServ := csrvim.NewCartService(categoryRepo)
//
//
//	cate := urepimp.NewUserGormRepo(nil)
//	cateServ := usrvimp.NewUserService(cate)
//
//	itemRepo := psrvim.NewItemService(nil)//NewMockCategoryRepo(nil)
//	v := psrvim.NewItemService(itemRepo)
//
//
//	CartHandler := NewCartHandler(tmpl,categoryServ,cateServ,v,nil)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/admin/carts", CartHandler.GetCarts)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	url := ts.URL
//
//	resp, err := tc.PostForm(url+"/deleteitemcart?id=1", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//
//}
//
//func TestGetSingleCart(t *testing.T) {
//
//	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))
//
//	categoryRepo := crepim.NewCartGormRepo(nil)
//	categoryServ := csrvim.NewCartService(categoryRepo)
//
//	CartHandler := NewCartHandler(tmpl,categoryServ,nil,nil,nil)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/admin/cart", CartHandler.GetSingleCart)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	sURL := ts.URL
//
//	resp, err := tc.PostForm(sURL+"/admin/cart?id=1", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//
//}
//
//func TestGetUserCart(t *testing.T) {
//
//	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))
//
//	categoryRepo := crepim.NewCartGormRepo(nil)
//	categoryServ := csrvim.NewCartService(categoryRepo)
//	cate := urepimp.NewUserGormRepo(nil)
//	cateServ := usrvimp.NewUserService(cate)
//
//	itemRepo := psrvim.NewItemService(nil)//NewMockCategoryRepo(nil)
//	v := psrvim.NewItemService(itemRepo)
//
//	CartHandler := NewCartHandler(tmpl,categoryServ,cateServ,v,nil)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/getusercart", CartHandler.GetUserCart)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	sURL := ts.URL
//
//	resp, err := tc.PostForm(sURL+"/getusercart?id=1", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//
//
//}
//
//func TestDeleteCart(t *testing.T) {
//
//	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))
//
//	categoryRepo := crepim.NewCartGormRepo(nil)
//	categoryServ := csrvim.NewCartService(categoryRepo)
//	cate := urepimp.NewUserGormRepo(nil)
//	cateServ := usrvimp.NewUserService(cate)
//
//	itemRepo := psrvim.NewItemService(nil)//NewMockCategoryRepo(nil)
//	v := psrvim.NewItemService(itemRepo)
//
//	CartHandler := NewCartHandler(tmpl,categoryServ,cateServ,v,nil)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/deleteitemcart", CartHandler.UpdateCart)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	sURL := ts.URL
//
//
//	resp, err := tc.PostForm(sURL+"/deleteitemcart?id=1", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//
//
//}
//
//
//func TestAddtoCart(t *testing.T) {
//
//	tmpl := template.Must(template.ParseGlob("../../../frontend/ui/templates/*"))
//
//	categoryRepo := crepim.NewCartGormRepo(nil)
//	categoryServ := csrvim.NewCartService(categoryRepo)
//	cate := urepimp.NewUserGormRepo(nil)
//	cateServ := usrvimp.NewUserService(cate)
//
//	itemRepo := psrvim.NewItemService(nil)//NewMockCategoryRepo(nil)
//	v := psrvim.NewItemService(itemRepo)
//
//	CartHandler := NewCartHandler(tmpl,categoryServ,cateServ,v,nil)
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/addtocart", CartHandler.AddtoCart)
//	ts := httptest.NewTLSServer(mux)
//	defer ts.Close()
//
//	tc := ts.Client()
//	sURL := ts.URL
//
//
//
//	resp, err := tc.PostForm(sURL+"/addtocart",nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//
//
//}
//
//
