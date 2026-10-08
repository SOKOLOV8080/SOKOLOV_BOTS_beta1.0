package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"SOKOLOV_BOTS/internal/database"
	"SOKOLOV_BOTS/internal/handlers"
	"SOKOLOV_BOTS/internal/models"

	"github.com/go-chi/chi/v5"
)

type PageData struct {
	User    models.User
	IsAuth  bool
	Content interface{}
}

func main() {
	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	r := chi.NewRouter()

	templates := template.Must(template.ParseGlob("templates/*.html"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, isAuth := handlers.GetCurrentUser(r)
		data := PageData{User: user, IsAuth: isAuth}
		templates.ExecuteTemplate(w, "index.html", data)
	})

	r.Get("/bots", func(w http.ResponseWriter, r *http.Request) {
		products, err := database.GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		user, isAuth := handlers.GetCurrentUser(r)
		data := PageData{User: user, IsAuth: isAuth, Content: products}
		templates.ExecuteTemplate(w, "bots.html", data)
	})

	r.Get("/bots/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		product, err := database.GetProductByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		user, isAuth := handlers.GetCurrentUser(r)
		data := PageData{User: user, IsAuth: isAuth, Content: product}
		templates.ExecuteTemplate(w, "product.html", data)
	})

	r.Get("/bots/{id}/download", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		product, err := database.GetProductByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if product.FilePath == "" {
			http.Error(w, "Файл не найден", 404)
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(product.FilePath)))
		http.ServeFile(w, r, product.FilePath)
	})

	r.Get("/sell", func(w http.ResponseWriter, r *http.Request) {
		user, isAuth := handlers.GetCurrentUser(r)
		if !isAuth {
			http.Redirect(w, r, "/login", 303)
			return
		}
		data := PageData{User: user, IsAuth: isAuth}
		templates.ExecuteTemplate(w, "sell.html", data)
	})

	r.Post("/sell", func(w http.ResponseWriter, r *http.Request) {
		user, isAuth := handlers.GetCurrentUser(r)
		if !isAuth {
			http.Redirect(w, r, "/login", 303)
			return
		}

		r.ParseMultipartForm(10 << 20)

		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

		file, header, err := r.FormFile("bot_file")
		if err != nil {
			http.Error(w, "Ошибка загрузки файла", 400)
			return
		}
		defer file.Close()

		filePath := fmt.Sprintf("uploads/%s", header.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", 500)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)

		p := models.Product{
			SellerID:    user.ID,
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Price:       price,
			GithubURL:   r.FormValue("github_url"),
			ContactInfo: r.FormValue("contact_info"),
			FilePath:    filePath,
		}

		err = database.CreateProduct(p)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/bots", 303)
	})

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "register.html", nil)
	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := database.CreateUser(email, password)
		if err != nil {
			data := map[string]string{"Error": "Пользователь с таким email уже существует"}
			templates.ExecuteTemplate(w, "register.html", data)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "user_email",
			Value: email,
			Path:  "/",
		})

		http.Redirect(w, r, "/", 303)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "login.html", nil)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := database.GetUserByEmail(email)
		if err != nil {
			data := map[string]string{"Error": "Пользователь не найден"}
			templates.ExecuteTemplate(w, "login.html", data)
			return
		}

		if !database.CheckPassword(user.PasswordHash, password) {
			data := map[string]string{"Error": "Неверный пароль"}
			templates.ExecuteTemplate(w, "login.html", data)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "user_email",
			Value: email,
			Path:  "/",
		})

		http.Redirect(w, r, "/", 303)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "user_email",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		http.Redirect(w, r, "/", 303)
	})
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		user, isAuth := handlers.GetCurrentUser(r)
		if !isAuth {
			http.Redirect(w, r, "/login", 303)
			return
		}

		products, err := database.GetProductsBySellerID(user.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		data := PageData{User: user, IsAuth: isAuth, Content: products}
		templates.ExecuteTemplate(w, "dashboard.html", data)
	})

	r.Post("/bots/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		user, isAuth := handlers.GetCurrentUser(r)
		if !isAuth {
			http.Redirect(w, r, "/login", 303)
			return
		}

		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		err := database.DeleteProduct(id, user.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/dashboard", 303)
	})
	log.Println("Server starting on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
