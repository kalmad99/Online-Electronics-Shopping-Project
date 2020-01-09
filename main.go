package main

import (
	"./entity"
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"./productlist/repository"
	"./productlist/service"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strconv"
	//uuid "github.com/satori/go.uuid"
)

var tmpl = template.Must(template.ParseGlob("delivery/web/templates/*.html"))
var productService *service.ProductService

func index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	products, err := productService.Products()
	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(w, "index.layout", products)
	if err != nil {
		panic(err)
	}
}

func seller(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "seller.index.layout", nil)
}

func indexProducts(w http.ResponseWriter, r *http.Request) {
	products, err := productService.Products()
	if err != nil {
		panic(err)
	}
	_ = tmpl.ExecuteTemplate(w, "seller.products.layout", products)
}

func sellerNewProducts(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		ctg := entity.Product{}
		ctg.Name = r.FormValue("name")
		ctg.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		ctg.ItemType = r.FormValue("type")
		ctg.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
		ctg.Description = r.FormValue("description")

		mf, fh, err := r.FormFile("catimg")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		ctg.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		err = productService.StoreProduct(ctg)

		if err != nil {
			_, _ = w.Write([]byte("Data Creation has failed or the item already exists"))
			panic(err)
		}

		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)

	} else {

		err := tmpl.ExecuteTemplate(w, "seller.product.new.layout", nil)

		if err != nil {
			panic(err)
		}

	}
}

func sellerUpdateProducts(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		cat, err := productService.Product(id)

		if err != nil {
			panic(err)
		}

		_ = tmpl.ExecuteTemplate(w, "seller.products.update.layout", cat)

	} else if r.Method == http.MethodPost {

		prod := entity.Product{}
		prod.ID, _ = strconv.Atoi(r.FormValue("id"))
		prod.Name = r.FormValue("name")
		prod.Description = r.FormValue("description")
		prod.Image = r.FormValue("image")

		mf, _, err := r.FormFile("catimg")

		if err != nil {
			panic(err)
		}

		defer mf.Close()

		writeFile(&mf, prod.Image)

		err = productService.UpdateProduct(prod)

		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
	}

}

func sellerDeleteProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		err = productService.DeleteProduct(id)

		if err != nil {
			panic(err)
		}

	}

	http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
}

func regist(w http.ResponseWriter, req *http.Request) {
	_ = tmpl.ExecuteTemplate(w, "Registrationform.html", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	_ = tmpl.ExecuteTemplate(w, "login.html", nil)
}
func Registration(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST" {
		http.Redirect(w, req, "/registrationpage", http.StatusSeeOther)
		return
	}
	usr := entity.User{}
	usr.Name = req.FormValue("name")
	usr.Email = req.FormValue("email")
	usr.Phone = req.FormValue("phone")
	usr.Password = req.FormValue("password")

	_ = tmpl.ExecuteTemplate(w, "update.html", usr)

	//err = userService.StoreUser(usr)
	//if err!=nil{
	//  http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	//  //panic(err.Error())
	//}
}
func writeFile(mf *multipart.File, fname string) {

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "delivery", "web", "assets", "img", fname)
	image, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	io.Copy(image, *mf)
}

func main() {
	dbDriver := "mysql"
	dbName := "golangtrialdb2"
	//dbName := "golangsession"
	dbUser := "root"
	dbPass := ""
	dbconn, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Println("Successfully connected to Mysql")
	//return dbconn

	//dbconn, err := sql.Open("postgres", "postgres://app_admin:P@$$w0rdD2@localhost/golangtrialdb?sslmode=disable")
	//
	//if err != nil {
	//	panic(err)
	//}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}

	proRepo := repository.NewPsqlProductRepository(dbconn)
	productService = service.NewProductService(proRepo)

	fs := http.FileServer(http.Dir("delivery/web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/seller", seller)
	http.HandleFunc("/seller/products", indexProducts)
	http.HandleFunc("/seller/products/new", sellerNewProducts)
	http.HandleFunc("/seller/products/update", sellerUpdateProducts)
	http.HandleFunc("/seller/products/delete", sellerDeleteProduct)
	http.HandleFunc("/registrationpage", regist)
	http.HandleFunc("/login", login)

	_ = http.ListenAndServe(":8181", nil)

}
