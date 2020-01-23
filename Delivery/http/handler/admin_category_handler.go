package handler

import (
	"github.com/kalmad99/Online-Electronics-Shopping-Project/allEntitiesAction/productpage"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/authFiles/csrfToken"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/entity"
	"github.com/kalmad99/Online-Electronics-Shopping-Project/frontend/form"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"


)

// AdminCategoryHandler handles category handler admin requests
type AdminCategoryHandler struct {
	tmpl        *template.Template
	categorySrv productpage.CategoryService
	csrfSignKey []byte
}

// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
func NewAdminCategoryHandler(t *template.Template, cs productpage.CategoryService, csKey []byte) *AdminCategoryHandler {
	return &AdminCategoryHandler{tmpl: t, categorySrv: cs, csrfSignKey: csKey}
}

// AdminCategories handle requests on route /admin/categories
func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request) {
	categories, errs := ach.categorySrv.Categories()
	if errs != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmplData := struct {
		Values     url.Values
		VErrors    form.ValidationErrors
		Categories []entity.Category
		CSRF       string
	}{
		Values:     nil,
		VErrors:    nil,
		Categories: categories,
		CSRF:       token,
	}
	ach.tmpl.ExecuteTemplate(w, "admin.categ.layout", tmplData)
}

// AdminCategoriesNew hanlde requests on route /admin/categories/new
func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		newCatForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		newCatForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newCatForm.Required("catname", "catdesc")
		newCatForm.MinLength("catdesc", 10)
		newCatForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !newCatForm.Valid() {
			ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
			return
		}
		mf, fh, err := r.FormFile("catimg")
		if err != nil {
			newCatForm.VErrors.Add("catimg", "File error")
			ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
			return
		}
		defer mf.Close()
		ctg := &entity.Category{
			Name:        r.FormValue("catname"),
			Description: r.FormValue("catdesc"),
			Image:       fh.Filename,
		}
		writeFile(&mf, fh.Filename)
		_, errs := ach.categorySrv.StoreCategory(ctg)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
	}
}

// AdminCategoriesUpdate handle requests on /admin/categories/update
func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		cat, errs := ach.categorySrv.Category(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		values := url.Values{}
		values.Add("catid", idRaw)
		values.Add("catname", cat.Name)
		values.Add("catdesc", cat.Description)
		values.Add("catimg", cat.Image)
		upCatForm := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			Category *entity.Category
			CSRF     string
		}{
			Values:   values,
			VErrors:  form.ValidationErrors{},
			Category: cat,
			CSRF:     token,
		}
		ach.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", upCatForm)
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
		updateCatForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		updateCatForm.Required("catname", "catdesc")
		updateCatForm.MinLength("catdesc", 10)
		updateCatForm.CSRF = token

		catID, err := strconv.Atoi(r.FormValue("catid"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		ctg := &entity.Category{
			ID:          uint(catID),
			Name:        r.FormValue("catname"),
			Description: r.FormValue("catdesc"),
			Image:       r.FormValue("imgname"),
		}
		mf, fh, err := r.FormFile("catimg")
		if err == nil {
			ctg.Image = fh.Filename
			err = writeFile(&mf, ctg.Image)
		}
		if mf != nil {
			defer mf.Close()
		}
		_, errs := ach.categorySrv.UpdateCategory(ctg)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
		return
	}
}

// AdminCategoriesDelete handle requests on route /admin/categories/delete
func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := ach.categorySrv.DeleteCategory(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

// AdminCategoriesDelete handle requests on route /admin/categories/delete
func (ach *AdminCategoryHandler) ItemsinCategories(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		cat, _ := ach.categorySrv.Category(uint(id))

		prds, errs := ach.categorySrv.ItemsInCategory(cat)

		if len(errs) > 0 {
			return
		}
		ach.tmpl.ExecuteTemplate(w, "index.layout", prds)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func writeFile(mf *multipart.File, fname string) error {
	wd, err := os.Getwd()
	log.Println("Working dir", wd)
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "../", "../", "frontend", "ui", "assets", "img", fname)
	image, err := os.Create(path)
	if err != nil {
		return err
	}
	defer image.Close()
	io.Copy(image, *mf)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
//package handler
//
//import (
//	"encoding/json"
//	"html/template"
//	"io"
//	"mime/multipart"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strconv"
//
//	"../../../entity"
//	"../../../productpage"
//)
//
//// AdminCategoryHandler handles category handler admin requests
//type AdminCategoryHandler struct {
//	tmpl        *template.Template
//	categorySrv productpage.CategoryService
//}
//
//// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
//func NewAdminCategoryHandler(t *template.Template, cs productpage.CategoryService) *AdminCategoryHandler {
//	return &AdminCategoryHandler{tmpl: t, categorySrv: cs}
//}
//
//// AdminCategories handle requests on route /admin/categories
//func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request) {
//	categories, errs := ach.categorySrv.Categories()
//
//	if len(errs) > 0 {
//		w.Header().Set("Content-type", "application/json")
//		http.Error(w, http.StatusText(http.StatusSeeOther), 303)
//		return
//	}
//	output, err := json.MarshalIndent(categories, "", "\t\t")
//
//	if err != nil{
//		w.Header().Set("Content-type", "application/json")
//		http.Error(w, http.StatusText(http.StatusSeeOther), 303)
//		return
//	}
//	_, err = w.Write(output)
//	if err != nil {
//		panic(err.Error())
//	}
//	//ach.tmpl.ExecuteTemplate(w, "admin.categ.layout", categories)
//}
//
//// AdminCategoriesNew hanlde requests on route /admin/categories/new
//func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodPost {
//
//		ctg := &entity.Category{}
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//
//		mf, fh, err := r.FormFile("catimg")
//		if err != nil {
//			panic(err)
//		}
//		defer mf.Close()
//
//		ctg.Image = fh.Filename
//
//		writeFile(&mf, fh.Filename)
//
//		_, errs := ach.categorySrv.StoreCategory(ctg)
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", nil)
//
//	}
//}
//
//// AdminCategoriesUpdate handle requests on /admin/categories/update
//func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
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
//		cat, errs := ach.categorySrv.Category(uint(id))
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", cat)
//
//	} else if r.Method == http.MethodPost {
//
//		ctg := &entity.Category{}
//		id, _ := strconv.Atoi(r.FormValue("id"))
//		ctg.ID = uint(id)
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//		ctg.Image = r.FormValue("image")
//
//		mf, _, err := r.FormFile("catimg")
//
//		if err != nil {
//			panic(err)
//		}
//
//		defer mf.Close()
//
//		writeFile(&mf, ctg.Image)
//
//		_, errs := ach.categorySrv.UpdateCategory(ctg)
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//	}
//
//}
//
//// AdminCategoriesDelete handle requests on route /admin/categories/delete
//func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request) {
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
//		_, errs := ach.categorySrv.DeleteCategory(uint(id))
//
//		if len(errs) > 0 {
//			panic(err)
//		}
//
//	}
//
//	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//}
//
//// AdminCategoriesDelete handle requests on route /admin/categories/delete
//func (ach *AdminCategoryHandler) ItemsinCategories(w http.ResponseWriter, r *http.Request) {
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
//		cat, _ := ach.categorySrv.Category(uint(id))
//
//		prds, errs := ach.categorySrv.ItemsInCategory(cat)
//
//		if len(errs) > 0 {
//			return
//		}
//		ach.tmpl.ExecuteTemplate(w, "index.layout", prds)
//	}
//
//	http.Redirect(w, r, "/", http.StatusSeeOther)
//}
//
//func writeFile(mf *multipart.File, fname string) {
//
//	wd, err := os.Getwd()
//
//	if err != nil {
//		panic(err)
//	}
//
//	path := filepath.Join(wd, "../", "../", "ui", "assets", "img", fname)
//	image, err := os.Create(path)
//
//	if err != nil {
//		panic(err)
//	}
//	defer image.Close()
//	io.Copy(image, *mf)
//}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
//AdminCategoryHandler handles category handler admin requests
//type AdminCategoryHandler struct {
//	tmpl        *template.Template
//	categorySrv productpage.CategoryService
//}
//
//// NewAdminCategoryHandler initializes and returns new AdminCateogryHandler
//func NewAdminCategoryHandler(t *template.Template, cs productpage.CategoryService) *AdminCategoryHandler {
//	return &AdminCategoryHandler{tmpl: t, categorySrv: cs}
//}
//
//// AdminCategories handle requests on route /admin/categories
//func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request) {
//	categories, errs := ach.categorySrv.Categories()
//	if errs != nil {
//		panic(errs)
//	}
//	ach.tmpl.ExecuteTemplate(w, "admin.categ.layout", categories)
//}
//
//// AdminCategoriesNew hanlde requests on route /admin/categories/new
//func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method == http.MethodPost {
//
//		ctg := entity.Category{}
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//
//		mf, fh, err := r.FormFile("catimg")
//		if err != nil {
//			panic(err)
//		}
//		defer mf.Close()
//
//		ctg.Image = fh.Filename
//
//		writeFile(&mf, fh.Filename)
//
//		_, err = ach.categorySrv.StoreCategory(ctg)
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", nil)
//
//	}
//}
//
//// AdminCategoriesUpdate handle requests on /admin/categories/update
//func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
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
//		cat, err := ach.categorySrv.Category(uint(id))
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		ach.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", cat)
//
//	} else if r.Method == http.MethodPost {
//
//		ctg := entity.Category{}
//		id, _ := strconv.Atoi(r.FormValue("id"))
//		ctg.ID = uint(id)
//		ctg.Name = r.FormValue("name")
//		ctg.Description = r.FormValue("description")
//		ctg.Image = r.FormValue("image")
//
//		mf, _, err := r.FormFile("catimg")
//
//		if err != nil {
//			panic(err)
//		}
//
//		defer mf.Close()
//
//		writeFile(&mf, ctg.Image)
//
//		err = ach.categorySrv.UpdateCategory(ctg)
//
//		if err != nil{
//			panic(err.Error())
//		}
//
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//
//	} else {
//		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//	}
//
//}
//
//// AdminCategoriesDelete handle requests on route /admin/categories/delete
//func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request) {
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
//		err = ach.categorySrv.DeleteCategory(uint(id))
//
//		if err != nil{
//			panic(err.Error())
//		}
//	}
//
//	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
//}
//
//func writeFile(mf *multipart.File, fname string) {
//
//	wd, err := os.Getwd()
//
//	if err != nil {
//		panic(err)
//	}
//
//	path := filepath.Join(wd, "../", "../", "ui", "assets", "img", fname)
//	image, err := os.Create(path)
//
//	if err != nil {
//		panic(err)
//	}
//	defer image.Close()
//	io.Copy(image, *mf)
//}
//
