package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/kalmad99/Online-Electronics-Shopping-Project/Delivery/http/handler"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	prepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/repository"
	psrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage/service"

	urepimp "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user/repository"
	usrvimp "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user/service"

	crepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/repository"
	csrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart/service"

	orepim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order/repository"
	osrvim "github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order/usecase"
)

func createTables(dbconn *gorm.DB) []error {
	errs := dbconn.CreateTable(&entity.User{}, &entity.Role{}, &entity.Session{}, &entity.Product{}, &entity.Bank{}, &entity.Category{}).GetErrors()
	//errs := dbconn.CreateTable(&entity.Bank{}).GetErrors()
	if errs != nil {
		return errs
	}
	return nil
}

func main() {

	csrfSignKey := []byte(csrfToken.GenerateRandomID(32))
	tmpl := template.Must(template.ParseGlob("../../frontend/ui/templates/*"))

	//dbconn, err := gorm.Open("postgres", "postgres://postgres:P@$$w0rdD2@localhost/restaurantdb?sslmode=disable")
	dbconn, err := gorm.Open("postgres", "postgres://postgres:password@localhost/restaurantdb4?sslmode=disable")

	createTables(dbconn)

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	sessionRepo := urepimp.NewSessionGormRepo(dbconn)
	sessionSrv := usrvimp.NewSessionService(sessionRepo)

	categoryRepo := prepim.NewCategoryGormRepo(dbconn)
	categoryServ := psrvim.NewCategoryService(categoryRepo)

	itemRepo := prepim.NewItemGormRepo(dbconn)
	itemServ := psrvim.NewItemService(itemRepo)

	userRepo := urepimp.NewUserGormRepo(dbconn)
	userServ := usrvimp.NewUserService(userRepo)

	cartRepo := crepim.NewCartGormRepo(dbconn)
	cartServ := csrvim.NewCartService(cartRepo)

	roleRepo := urepimp.NewRoleGormRepo(dbconn)
	roleServ := usrvimp.NewRoleService(roleRepo)

	orderRepo := orepim.NewOrderGormRepo(dbconn)
	orderServ := osrvim.NewOrderService(orderRepo)

	ach := handler.NewAdminCategoryHandler(tmpl, categoryServ, csrfSignKey)
	oh := handler.NewOrderHandler(tmpl, orderServ, userServ, itemServ, csrfSignKey)
	sph := handler.NewSellerProductHandler(tmpl, itemServ, csrfSignKey)
	mh := handler.NewMenuHandler(tmpl, itemServ, csrfSignKey)
	arh := handler.NewAdminRoleHandler(roleServ)

	sess := ConfigSessions()
	uh := handler.NewUserHandler(tmpl, userServ, sessionSrv, roleServ, sess, csrfSignKey)
	ch := handler.NewCartHandler(tmpl, cartServ, userServ, sessionSrv, roleServ, sess, itemServ, csrfSignKey)

	fs := http.FileServer(http.Dir("../../frontend/ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mh.Index)
	http.HandleFunc("/about", mh.About)
	http.HandleFunc("/contact", mh.Contact)
	http.HandleFunc("/menu", mh.Menu)
	http.Handle("/admin", uh.Authenticated(uh.Authorized(http.HandlerFunc(mh.Admin))))
	http.HandleFunc("/Loginpage", mh.LoginPage)
	http.HandleFunc("/Registpage", mh.RegistPage)

	http.Handle("/admin/users", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.Users))))
	http.Handle("/admin/categories", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategories))))
	http.Handle("/admin/categories/new", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesNew))))
	http.Handle("/admin/categories/update", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesUpdate))))
	http.Handle("/admin/categories/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesDelete))))
	//http.Handle("admin/category", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.ItemsinCategories))))

	http.Handle("/admin/roles/new", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.PostRole))))
	http.Handle("/admin/roles", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.GetRoles))))
	http.Handle("/admin/role", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.GetSingleRole))))
	http.Handle("/admin/roles/update", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.PutRole))))
	http.Handle("/admin/roles/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.DeleteRole))))

	http.HandleFunc("/category", ach.ItemsinCategories)

	http.Handle("/seller/products", uh.Authenticated(http.HandlerFunc(sph.SellerProducts)))
	//http.Handle("/seller/products",uh.Authenticated(http.HandlerFunc(sph.SellerProducts)))
	http.Handle("/seller/products/new", uh.Authenticated(http.HandlerFunc(sph.SellerProductsNew)))
	http.Handle("/seller/products/update", uh.Authenticated(http.HandlerFunc(sph.SellerProductsUpdate)))
	http.Handle("/seller/products/delete", uh.Authenticated(http.HandlerFunc(sph.SellerProductsDelete)))
	http.HandleFunc("/searchProducts", sph.SearchProducts)
	http.HandleFunc("/detail", sph.ProductDetail)
	http.HandleFunc("/rate", sph.Rating)

	http.Handle("/admin/carts", uh.Authenticated(uh.Authorized(http.HandlerFunc(ch.GetCarts))))
	http.Handle("/admin/cart", uh.Authenticated(uh.Authorized(http.HandlerFunc(ch.GetSingleCart))))
	
	http.Handle("/admin/orders", uh.Authenticated(uh.Authorized(http.HandlerFunc(oh.Orders))))
	http.Handle("/admin/order", uh.Authenticated(uh.Authorized(http.HandlerFunc(oh.GetUserOrder))))
	http.Handle("/admin/order/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(oh.OrderDelete))))

	http.Handle("/getusercart", uh.Authenticated(http.HandlerFunc(ch.GetUserCart)))
	//http.HandleFunc("/getusercart", ch.GetUserCart)
	http.Handle("/deleteitemcart", uh.Authenticated(http.HandlerFunc(ch.UpdateCart)))
	//http.HandleFunc("/deleteitemcart", ch.UpdateCart)

	//http.HandleFunc("/addtocart", ch.AddtoCart)
	http.Handle("/addtocart", uh.Authenticated(http.HandlerFunc(ch.AddtoCart)))


	http.HandleFunc("/users/success", uh.UsersUpdate)
	http.HandleFunc("/user/update", uh.UsersUpdate)
	//http.HandleFunc("/admin/users",userHandler.Users)
	http.HandleFunc("/registrationprocess1", uh.Signup)
	http.HandleFunc("/Registration2", uh.Registration2)
	http.HandleFunc("/user/delete", uh.UsersDelete)
	//http.HandleFunc("/user/changepass", uh.UserChangePassword)
	http.HandleFunc("/login", uh.Login)
	http.HandleFunc("/logout", uh.Logout)
	//http.HandleFunc("/cart", uh.Cart)
	http.ListenAndServe(":8080", nil)
}

func ConfigSessions() *entity.Session {
	tokenExpires := time.Now().Add(time.Minute * 30).Unix()
	sessionID := csrfToken.GenerateRandomID(32)
	signingString, err := csrfToken.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	signingKey := []byte(signingString)

	return &entity.Session{
		Expires:    tokenExpires,
		SigningKey: signingKey,
		UUID:       sessionID,
	}
}
