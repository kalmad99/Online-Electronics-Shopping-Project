package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"

	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/permission"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/session"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"golang.org/x/crypto/bcrypt"
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

var cid string
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
			// log.Println("Got here 65")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := uh.userService.UserRoles(uh.loggedInUser)
		if len(errs) > 0 {
			// log.Println("Got here 71")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
			if !permitted {
				// log.Println("Got here 79")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		if r.Method == http.MethodPost {
			ok, err := csrfToken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
			if !ok || (err != nil) {
				// log.Println("Got here 87")
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
			// UserID  string
			CSRF string
		}{
			Values:  nil,
			VErrors: nil,
			// UserID:  "",
			CSRF: token,
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
			loginForm.VErrors.Add("generic", "Your Email Address and/or Password is Wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(r.FormValue("password")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			loginForm.VErrors.Add("generic", "Your Email Address and/or Password is Wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}
		uh.loggedInUser = usr
		claims := csrfToken.Claims(usr.Email, uh.userSess.Expires)
		session.Create(claims, uh.userSess.UUID, uh.userSess.SigningKey, w)
		newSess, errs := uh.sessionService.StoreSession(uh.userSess)
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to Store Session")
			uh.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
			return
		}
		// usrid := fmt.Sprint(usr.ID)
		// newloginForm := struct{
		// 	loginForm  form.Input
		// 	userID string
		// }{
		// 	loginForm: loginForm,
		// 	userID: usrid,
		// }
		uh.userSess = newSess
		roles, _ := uh.userService.UserRoles(usr)
		if uh.checkAdmin(roles) {
			// uh.tmpl.ExecuteTemplate(w, "login.html", newloginForm)
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		cid = fmt.Sprint(usr.ID)
		link := "/?userid=" + fmt.Sprint(usr.ID)
		http.Redirect(w, r, link, http.StatusSeeOther)
		// http.Redirect(w, r, "/", http.StatusSeeOther)
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
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.Required("name", "email", "password", "confpass")
		signUpForm.MatchesPattern("email", form.EmailRX)
		signUpForm.MatchesPhonePattern("phone", form.PhoneRX)
		signUpForm.MinLength("password", 8)
		signUpForm.PasswordMatches("password", "confpass")
		signUpForm.CSRF = token
		//If there are any errors, redisplay the signup form.
		if !signUpForm.Valid() {
			err = uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		//log.Println("Got here 206")

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			signUpForm.VErrors.Add("phone", "Phone Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		//log.Println("Got here 217")
		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			signUpForm.VErrors.Add("email", "Email Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		log.Println("Got here 238")
		role, errs := uh.userRole.RoleByName("USER")

		if len(errs) > 0 {
			signUpForm.VErrors.Add("role", "could not assign role to the user")
			err := uh.tmpl.ExecuteTemplate(w, "Registrationform.layout", signUpForm)
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
	//log.Println("Logged in ", c.Value)
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

func (uh *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, errs := uh.userService.Users()
	if errs != nil {
		panic(errs)
	}
	uh.tmpl.ExecuteTemplate(w, "admin.users.layout", users)
}

func (uh *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet{
		idraw := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idraw)
		usr, errs := uh.userService.User(uint(id))
		if errs != nil {
			panic(errs)
		}
		userProf := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
			User    *entity.User
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
			User:    usr,
		}

		uh.tmpl.ExecuteTemplate(w, "user.index.layout", userProf)
	}
	if r.Method == http.MethodPost{
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
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
		err = uh.tmpl.ExecuteTemplate(w, "user.update.html", upAccForm)
		if err != nil {
			panic(err)
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
		http.Redirect(w, r, "/userprof?id=" + cid, http.StatusSeeOther)
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

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
