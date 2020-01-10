package http

import (
	"../../entity"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"

	"../../productpage/repository"
	"../../productpage/service"
	"../../users/urepository"
	"../../users/uservice"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	//uuid "github.com/satori/go.uuid"
	"../http/handler"
)

var name, email, phone, pass string
var id int

var tmpl = template.Must(template.ParseGlob("delivery/web/templates/*.html"))
var productService *service.ProductService
var userService *uservice.UserService

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

func regist(w http.ResponseWriter, req *http.Request) {
	_ = tmpl.ExecuteTemplate(w, "Registrationform.html", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	_ = tmpl.ExecuteTemplate(w, "login.html", nil)
}

func Registration(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
		return
	}
	usr := entity.User{}
	usr.Name = req.FormValue("name")
	usr.Email = req.FormValue("email")
	usr.Phone = req.FormValue("phone")
	usr.Password = req.FormValue("password")

	name = usr.Name
	email = usr.Email
	phone = usr.Phone
	pass = usr.Password

	hostURL := "smtp.gmail.com"
	hostPort := "587"
	emailSender := "kalemesfin12go@gmail.com"
	password := "qnzfgwbnaxykglvu"
	emailReceiver := usr.Email

	emailAuth := smtp.PlainAuth(
		"",
		emailSender,
		password,
		hostURL,
	)

	msg := []byte("To: " + emailReceiver + "\r\n" +
		"Subject: " + "Hello " + usr.Name + "\r\n" +
		"This is your OTP. 123456789")

	err := smtp.SendMail(
		hostURL+":"+hostPort,
		emailAuth,
		emailSender,
		[]string{emailReceiver},
		msg,
	)

	if err != nil {
		fmt.Print("Error: ", err)
	}
	fmt.Print("Email Sent")

	//err = userService.StoreUser(usr)
	//if err!=nil{
	//  panic(err.Error())
	//}

	//_ = tmpl.ExecuteTemplate(w, "Registrationformpart2.html", info)
	_ = tmpl.ExecuteTemplate(w, "registotp.html", usr)

	//err = userService.StoreUser(usr)
	//if err!=nil{
	//  http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	//  //panic(err.Error())
	//}
}
func Registration2(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
		return
	}

	otp := req.FormValue("otpfield")

	usrinfo := entity.User{uint(id), name, email, phone, pass}

	if otp == "123456789" {
		//_ = tpl.ExecuteTemplate(w, "update.html", usrinfo)
		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
		err := userService.StoreUser(usrinfo)
		if err != nil {
			http.Redirect(w, req, "/Registration2", http.StatusSeeOther)
			//panic(err.Error())
		}
	} else {
		fmt.Print("Wrong otp")
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	}
	http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	return
}
func Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
		return
	}
	email := req.FormValue("email")
	password := req.FormValue("password")

	log.Println(email)
	usr, err := userService.User(email)

	//log.Println(usr.Name)
	//log.Println(usr.Email)
	//log.Println(usr.Phone)
	//log.Println(usr.Password)

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		log.Println("Username or Password is incorrect")
		http.Redirect(w, req, "/Loginpage", 301)
		return
	}
	err = tmpl.ExecuteTemplate(w, "update.html", usr)
	if err != nil {
		panic(err.Error())
	}

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
	fmt.Println("Successfully connected to Mysql")
	//return dbconn

	//dbconn, err := sql.Open("postgres", "postgres://postgres:Dadusha99@localhost/restaurantdb2?sslmode=disable")
	//
	//if err != nil {
	//	panic(err)
	//}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}
	productRepo := repository.NewPsqlProductRepository(dbconn)
	productServ := service.NewProductService(productRepo)
	sellerProHandler := handler.NewSellerProductHandler(tmpl, productServ)

	//proRepo := repository.NewPsqlProductRepository(dbconn)
	//productService = service.NewProductService(proRepo)

	usrRepo := urepository.NewPsqlUserRepository(dbconn)
	userService = uservice.NewUserService(usrRepo)

	fs := http.FileServer(http.Dir("delivery/web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/seller", seller)
	http.HandleFunc("/seller/products", sellerProHandler.SellerProducts)
	http.HandleFunc("/seller/products/new", sellerProHandler.SellerProductsNew)
	http.HandleFunc("/seller/products/update", sellerProHandler.SellerProductsUpdate)
	http.HandleFunc("/seller/products/delete", sellerProHandler.SellerProductsDelete)
	http.HandleFunc("/registrationpage", regist)
	http.HandleFunc("/Registration", Registration)
	http.HandleFunc("/login", Login)

	_ = http.ListenAndServe(":8181", nil)

}
