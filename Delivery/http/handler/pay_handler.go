package handler

import (
	"fmt"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/bank"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/cart"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/order"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/user"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
)

type PayHandler struct {
	tmpl        *template.Template
	bankSrv     bank.PayService
	userSrv     user.UserService
	orderSrv    order.OrderService
	cartSrv     cart.CartService
	csrfSignKey []byte
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewPayHandler(t *template.Template, bs bank.PayService, us user.UserService, os order.OrderService,
	carts cart.CartService, cs []byte) *PayHandler {
	return &PayHandler{tmpl: t, bankSrv: bs, userSrv: us, orderSrv: os, cartSrv: carts, csrfSignKey: cs}
}

func (ph *PayHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(ph.csrfSignKey)
	ordr := entity.Order{}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idRaw)
		ordr.ID = uint(id)
		log.Println("Get order ", ordr.ID)

		newPayForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			Order   entity.Order
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			Order:   ordr,
			CSRF:    token,
		}
		err := ph.tmpl.ExecuteTemplate(w, "payment.html", newPayForm)
		if err != nil {
			panic(err)
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		newPayForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newPayForm.Required("bank")
		newPayForm.CSRF = token

		log.Println("This is balance bank", r.FormValue("bank"))
		bExists := ph.bankSrv.BankExists(r.FormValue("bank"))
		log.Println("Exists", bExists)
		if !bExists {
			oid := entity.User{}
			i, _ := strconv.Atoi(idr)
			oid.ID = uint(i)
			ordr, _ := ph.orderSrv.CustomerOrders(&oid)
			err := ph.tmpl.ExecuteTemplate(w, "payerror.html", ordr.ID)
			if err != nil {
				panic(err)
			}
			return
		} else {
			bankacc := r.FormValue("bank")
			userid := r.FormValue("order")
			uid, _ := strconv.Atoi(userid)
			order, errs := ph.orderSrv.Order(uint(uid))
			userord, _ := ph.userSrv.User(order.UserID)
			if len(errs) > 0 {
				panic(errs)
			}
			log.Println("Amount", order.Total)
			log.Println("user", order.UserID)
			log.Println("user email", userord.Email)
			_, errs = ph.bankSrv.MakePayment(bankacc, order.Total)
			if len(errs) > 0 {
				panic(errs)
			}
			hostURL := "smtp.gmail.com"
			hostPort := "587"
			emailSender := "kalemesfin12go@gmail.com"
			password := "qnzfgwbnaxykglvu"
			emailReceiver := userord.Email

			emailAuth := smtp.PlainAuth(
				"",
				emailSender,
				password,
				hostURL,
			)

			tid := RandStringBytes(10)
			msg := []byte("To: " + emailReceiver + "\r\n" +
				"Subject: " + "E&E Online Shopping" + "\r\n" +
				"Dear " + userord.Name + "\r\n" +
				"Your Order was successfully made. Your Transaction ID is " + tid + "\r\n" +
				"Thank you for using our Website!!!")

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

			final := struct {
				User          entity.User
				TransactionID string
			}{
				User:          *userord,
				TransactionID: tid,
			}
			_, errs = ph.orderSrv.DeleteOrder(uint(uid))
			if len(errs) > 0 {
				panic(errs)
			}
			_, errs = ph.cartSrv.DeleteCart(userord)
			ph.tmpl.ExecuteTemplate(w, "pay.success.html", final)
		}
	}
}
