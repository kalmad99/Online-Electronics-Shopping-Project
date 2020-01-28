package handler

import (
	"context"
	"fmt"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/permission"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/session"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
)

// UserHandler handler handles user related requests
type UserHandler struct {
	tmpl           *template.Template
	userService    user.UserService
	sessionService user.SessionService
	userSess       *entity.Session
	loggedInUser   *entity.User
	userRole       user.RoleService
	csrfSignKey    []byte
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")
var name, email, phone, pass string
var id, roleid uint

// NewUserHandler returns new UserHandler object
func NewUserHandler(t *template.Template, usrServ user.UserService,
	sessServ user.SessionService, uRole user.RoleService,
	usrSess *entity.Session, csKey []byte) *UserHandler {
	return &UserHandler{tmpl: t, userService: usrServ, sessionService: sessServ,
		userRole: uRole, userSess: usrSess, csrfSignKey: csKey}
}

// Authenticated checks if a user is authenticated to access a given route
func (uh *UserHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := uh.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, uh.userSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Authorized checks if a user has proper authority to access a give route
func (uh *UserHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uh.loggedInUser == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := uh.userService.UserRoles(uh.loggedInUser)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
			if !permitted {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		if r.Method == http.MethodPost {
			ok, err := csrfToken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
			if !ok || (err != nil) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Login hanldes the GET/POST /login requests
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		loginForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
		return
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		usr, errs := uh.userService.Login(r.FormValue("email"))
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(r.FormValue("password")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			loginForm.VErrors.Add("generic", "Your email address or password is wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}

		uh.loggedInUser = usr
		claims := csrfToken.Claims(usr.Email, uh.userSess.Expires)
		session.Create(claims, uh.userSess.UUID, uh.userSess.SigningKey, w)
		newSess, errs := uh.sessionService.StoreSession(uh.userSess)
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to store session")
			uh.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
			return
		}
		uh.userSess = newSess
		roles, _ := uh.userService.UserRoles(usr)
		if uh.checkAdmin(roles) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		//err = uh.tmpl.ExecuteTemplate(w, "index.layout", usr)
		//if err != nil {
		//	panic(err)
		//}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout hanldes the POST /logout requests
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	session.Remove(userSess.UUID, w)
	uh.sessionService.DeleteSession(userSess.UUID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Signup hanldes the GET/POST /registrationprocess1 requests
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		signUpForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		err := uh.tmpl.ExecuteTemplate(w, "Registrationform.html", signUpForm)
		if err != nil {
			panic(err.Error())
		}
		return
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate the form contents
		//log.Println("Got here 190")
		singnUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("name", "email", "password", "confpass")
		singnUpForm.MatchesPattern("email", form.EmailRX)
		singnUpForm.MatchesPhonePattern("phone", form.PhoneRX)
		singnUpForm.MinLength("password", 8)
		singnUpForm.PasswordMatches("password", "confpass")
		singnUpForm.CSRF = token
		//If there are any errors, redisplay the signup form.
		if !singnUpForm.Valid() {
			err = uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", singnUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		//log.Println("Got here 206")

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			singnUpForm.VErrors.Add("phone", "Phone Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", singnUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		//log.Println("Got here 217")
		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			singnUpForm.VErrors.Add("email", "Email Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", singnUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		log.Println("Got here 227")
		//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 12)
		//if err != nil {
		//	singnUpForm.VErrors.Add("password", "Password Could not be stored")
		//	err := uh.tmpl.ExecuteTemplate(w, "signup.layout", singnUpForm)
		//	if err != nil {
		//		panic(err.Error())
		//	}
		//	return
		//}

		log.Println("Got here 238")
		role, errs := uh.userRole.RoleByName("USER")

		if len(errs) > 0 {
			singnUpForm.VErrors.Add("role", "could not assign role to the user")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", singnUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		log.Println("Got here 245")
		user := &entity.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: r.FormValue("password"),
			RoleID:   role.ID,
		}

		name = user.Name
		email = user.Email
		phone = user.Phone
		pass = user.Password
		id = user.ID
		roleid = user.RoleID

		hostURL := "smtp.gmail.com"
		hostPort := "587"
		emailSender := "kalemesfin12go@gmail.com"
		password := "qnzfgwbnaxykglvu"
		emailReceiver := user.Email

		emailAuth := smtp.PlainAuth(
			"",
			emailSender,
			password,
			hostURL,
		)

		msg := []byte("To: " + emailReceiver + "\r\n" +
			"Subject: " + "Hello " + user.Name + "\r\n" +
			"This is your OTP. 123456789")

		err = smtp.SendMail(
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

		_ = uh.tmpl.ExecuteTemplate(w, "Registrationformpart2.html", user)
	}
}

//Second stage registration /Registration2
func (uh *UserHandler) Registration2(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
		return
	}
	otp := req.FormValue("otpfield")

	usrinfo := &entity.User{ID: id, Name: name, Email: email, Phone: phone, Password: pass, RoleID: roleid}

	if otp == "123456789" {
		//_ = tpl.ExecuteTemplate(w, "update.html", usrinfo)
		_, err := uh.userService.StoreUser(usrinfo)
		if err != nil {
			http.Redirect(w, req, "/Registration2", http.StatusSeeOther)
			//panic(err.Error())
		}
		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
	} else {
		fmt.Print("Wrong otp")
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	}
	http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	return

}
func (uh *UserHandler) loggedIn(r *http.Request) bool {
	if uh.userSess == nil {
		return false
	}
	userSess := uh.userSess
	c, err := r.Cookie(userSess.UUID)
	if err != nil {
		return false
	}
	ok, err := session.Valid(c.Value, userSess.SigningKey)
	if !ok || (err != nil) {
		return false
	}
	return true
}

func (uh *UserHandler) checkAdmin(rs []entity.Role) bool {
	for _, r := range rs {
		if strings.ToUpper(r.Name) == strings.ToUpper("Admin") {
			return true
		}
	}
	return false
}

// AdminCategoryHandler handles category handler admin requests
//type UserHandler struct {
//	tmpl        *template.Template
//	userSrv user.UserService
//}
//
//// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
//func NewUserHandler(t *template.Template, us user.UserService) *UserHandler {
//	return &UserHandler{tmpl: t, userSrv: us}
//}
//
// Users handle requests on route /admin/users
func (uh *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, errs := uh.userService.Users()
	if errs != nil {
		panic(errs)
	}
	uh.tmpl.ExecuteTemplate(w, "admin.user.layout", users)
}

//
//// AdminCategoriesNew hanlde requests on route /admin/categories/new
//func (uh *UserHandler) UserNew(w http.ResponseWriter, req *http.Request) {
//
//	if req.Method != "POST" {
//		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
//		return
//	}
//
//	usr := entity.User{}
//	usr.Name = req.FormValue("name")
//	usr.Email = req.FormValue("email")
//	usr.Phone = req.FormValue("phone")
//	usr.Password = req.FormValue("password")
//
//	name = usr.Name
//	email = usr.Email
//	phone = usr.Phone
//	pass = usr.Password
//
//	hostURL := "smtp.gmail.com"
//	hostPort := "587"
//	emailSender := "kalemesfin12go@gmail.com"
//	password := "qnzfgwbnaxykglvu"
//	emailReceiver := usr.Email
//
//	emailAuth := smtp.PlainAuth(
//		"",
//		emailSender,
//		password,
//		hostURL,
//	)
//
//	msg := []byte("To: " + emailReceiver + "\r\n" +
//		"Subject: " + "Hello " + usr.Name + "\r\n" +
//		"This is your OTP. 123456789")
//
//	err:=	smtp.SendMail(
//		hostURL + ":" + hostPort,
//		emailAuth,
//		emailSender,
//		[]string{emailReceiver},
//		msg,
//	)
//
//	if err != nil{
//		fmt.Print("Error: ", err)
//	}
//	fmt.Print("Email Sent")
//
//	//err = userService.StoreUser(usr)
//	//if err!=nil{
//	//	panic(err.Error())
//	//}
//
//	//_ = tmpl.ExecuteTemplate(w, "Registrationformpart2.html", info)
//	_ = uh.tmpl.ExecuteTemplate(w, "registotp.html", usr)
//
//	//err = userService.StoreUser(usr)
//	//if err!=nil{
//	//	http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
//	//	//panic(err.Error())
//	//}
//}
//func (uh *UserHandler) Registration2(w http.ResponseWriter, req *http.Request) {
//	if req.Method != "POST"{
//		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
//		return
//	}
//	otp := req.FormValue("otpfield")
//
//	usrinfo := &entity.User{ID:id, Name:name, Email:email, Phone:phone, Password:pass}
//
//	if otp == "123456789" {
//		//_ = tpl.ExecuteTemplate(w, "update.html", usrinfo)
//		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
//		_, err := uh.userSrv.StoreUser(usrinfo)
//		if err!=nil{
//			http.Redirect(w, req, "/Registration2", http.StatusSeeOther)
//			//panic(err.Error())
//		}
//	} else{
//		fmt.Print("Wrong otp")
//		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
//	}
//	http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
//	return
//}
//
//func (uh *UserHandler) UserLogin(w http.ResponseWriter, req *http.Request) {
//	if req.Method != "POST" {
//		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
//		return
//	}
//	email := req.FormValue("email")
//	password := req.FormValue("password")
//
//	log.Println(email)
//	usr, errs := uh.userSrv.Login(email)
//
//	if len(errs) > 0{
//		panic(errs)
//	}
//
//	//log.Println(usr.Name)
//	//log.Println(usr.Email)
//	//log.Println(usr.Phone)
//	//log.Println(usr.Password)
//
//	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
//	if err != nil {
//		log.Println("Username or Password is incorrect")
//		http.Redirect(w, req, "/Loginpage", 301)
//		return
//	}
//	err = uh.tmpl.ExecuteTemplate(w, "user.index.layout", usr)
//	if err != nil {
//		panic(err.Error())
//	}
//}
//
//// AdminCategoriesUpdate handle requests on /admin/categories/update
//func (uh *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		usr, errs := uh.userSrv.User(uint(id))
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		uh.tmpl.ExecuteTemplate(w, "user.update.layout", usr)
//
//	} else if r.Method == http.MethodPost {
//
//		usr := &entity.User{}
//		id, _ := strconv.Atoi(r.FormValue("id"))
//		usr.ID = uint(id)
//		usr.Name = r.FormValue("name")
//		usr.Email = r.FormValue("email")
//		usr.Phone = r.FormValue("phone")
//
//		_, errs := uh.userSrv.UpdateUser(usr)
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//		http.Redirect(w, r, "/Loginpage", http.StatusSeeOther)
//	}
//}
//
//

// UsersUpdate handles GET/POST /users/update?id={id} request
func (uh *UserHandler) UsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user, errs := uh.userService.User(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		values := url.Values{}
		values.Add("userid", idRaw)
		values.Add("name", user.Name)
		values.Add("email", user.Email)
		values.Add("phone", user.Phone)

		upAccForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			User    *entity.User
			CSRF    string
		}{
			Values:  values,
			VErrors: form.ValidationErrors{},
			User:    user,
			CSRF:    token,
		}
		uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
		return
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		upAccForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		upAccForm.Required("name", "email", "phone")
		upAccForm.MatchesPattern("email", form.EmailRX)
		upAccForm.MatchesPhonePattern("phone", form.PhoneRX)
		upAccForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !upAccForm.Valid() {
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}
		userID := r.FormValue("userid")
		uid, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user, errs := uh.userService.User(uint(uid))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if (user.Email != r.FormValue("email")) && eExists {
			upAccForm.VErrors.Add("email", "Email Already Exists")
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))

		if (user.Phone != r.FormValue("phone")) && pExists {
			upAccForm.VErrors.Add("phone", "Phone Already Exists")
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}
		if err != nil {
			upAccForm.VErrors.Add("role", "could not retrieve role id")
			uh.tmpl.ExecuteTemplate(w, "admin.user.update.layout", upAccForm)
			return
		}
		usr := &entity.User{
			ID:       user.ID,
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: user.Password,
			RoleID:   user.RoleID,
		}
		_, errs = uh.userService.UpdateUser(usr)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

// UsersDelete handles Delete /users/delete?id={id} request
func (uh *UserHandler) UsersDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := uh.userService.DeleteUser(uint(id))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
		session.Remove(userSess.UUID, w)
		uh.sessionService.DeleteSession(userSess.UUID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AdminUsersDelete handles Delete /admin/users/delete?id={id} request
func (uh *UserHandler) AdminUsersDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := uh.userService.DeleteUser(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

//// AdminCategoriesDelete handle requests on route /admin/categories/delete
//func (uh *UserHandler) UserDelete(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		_, errs := uh.userSrv.DeleteUser(uint(id))
//
//		if len(errs) > 0 {
//			panic(err)
//		}
//
//	}
//
//	http.Redirect(w, r, "/users", http.StatusSeeOther)
//}
//
//func (uh *UserHandler) UserChangePassword(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		usr, errs := uh.userSrv.User(uint(id))
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		uh.tmpl.ExecuteTemplate(w, "changepass.layout", usr)
//
//	}
//	usr := &entity.User{}
//	id, _ := strconv.Atoi(r.FormValue("id"))
//	//usr.ID = uint(id)
//	usr.ID = uint(id)
//	user, err1 := uh.userSrv.User(usr.ID)
//	if len(err1) > 0{
//		panic(err1)
//	}
//
//	log.Println(usr.ID)
//	usr.Password = r.FormValue("password")
//	var confp = r.FormValue("confpass")
//	var oldp = r.FormValue("oldpass")
//
//	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldp))
//	if err != nil {
//		log.Println("This is not your old password")
//		http.Redirect(w, r, "/users", 301)
//		return
//	}
//	if usr.Password == confp {
//
//		hashedpass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
//		if err != nil {
//			//http.Error(w, "Server error, unable to create your account.", 500)
//			//return errors.New("server error, unable to create your account")
//			panic(err.Error())
//		}
//		usr.Password=string(hashedpass)
//
//		_, errs := uh.userSrv.ChangePassword(usr)
//		log.Println("Succesfully changed")
//		http.Redirect(w, r, "/Loginpage", 303)
//		if len(errs) > 0{
//			panic(errs)
//		}
//	}else{
//		log.Println("Passwords don't match")
//		uh.tmpl.ExecuteTemplate(w, "user.update.layout", nil)
//	}
//
//
//
//	//usr := &entity.User{}
//	//id, _ := strconv.Atoi(r.FormValue("id"))
//	//usr.ID = uint(id)
//	//usr.Name = r.FormValue("name")
//	//usr.Email = r.FormValue("email")
//	//usr.Phone = r.FormValue("phone")
//	//usr.Password = r.FormValue("newpass")
//	//var confp = r.FormValue("confnewpass")
//	//
//	//if usr.Password == confp{
//	//	hashedpass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
//	//	if err != nil {
//	//		panic(err.Error())
//	//	}
//	//	usr.Password=string(hashedpass)
//	//
//	//	_, errs := uh.userSrv.ChangePassword(usr)
//	//
//	//	if len(errs) > 0 {
//	//		panic(errs)
//	//	}
//	//}else{
//	//	errors.New("The passwords you entered dont match")
//	//}
//	//http.Redirect(w, r, "/users", http.StatusSeeOther)
//}
