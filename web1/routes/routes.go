package routes

import (
	"github.com/gorilla/mux"
	"github.com/iitheogotut/web/web1/middleware"
	"github.com/iitheogotut/web/web1/models"
	"github.com/iitheogotut/web/web1/sessions"
	"github.com/iitheogotut/web/web1/utils"
	"log"
	"net/http"
)

func NewRouter() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	router.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	router.HandleFunc("/login", loginGetHandler).Methods("GET")
	router.HandleFunc("/logout", logoutGetHandler).Methods("GET")
	router.HandleFunc("/login", loginPostHandler).Methods("POST")
	router.HandleFunc("/register", registerGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/{username}", middleware.AuthRequired(userGetHandler)).Methods("GET")
	return router
}


func indexGetHandler(w http.ResponseWriter, r *http.Request){
	updates, err := models.GetAllUpdates()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct{
		Title string
		Updates []*models.Update
		DisplayForm bool
	}{
		Title: "All updates",
		Updates: updates,
		DisplayForm:true,
	})

}

func indexPostHandler(w http.ResponseWriter, r *http.Request){
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}
	body := r.PostForm.Get("update")
	err = models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w,r,"/",302)
}

func userGetHandler(w http.ResponseWriter, r *http.Request){
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	currentUserId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	updates, err := models.GetUpdates(userId)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct{
		Title string
		Updates []*models.Update
		DisplayForm bool
	}{
		Title: username,
		Updates: updates,
		DisplayForm: currentUserId == userId,
	})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request){
	utils.ExecuteTemplate(w, "login.html", nil)

}

func logoutGetHandler(w http.ResponseWriter, r *http.Request){
	session, _ := sessions.Store.Get(r, "session")
	delete(session.Values, "user_id")
	err := session.Save(r, w)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)

	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "invalid login")
		default:
			utils.InternalServerError(w)
		}
		return
	}

	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}

	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/", 302)

}

func registerGetHandler(w http.ResponseWriter, r *http.Request){
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request){
	_ = r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken{
		utils.ExecuteTemplate(w, "register.html", "username taken")
		return
	} else if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/login", 302)

}
