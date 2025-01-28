package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

// User структура для хранения пользователей
type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

var userCollection *mongo.Collection

// registerHandler обработчик для регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("static/register.html")
		if err != nil {
			http.Error(w, "Error loading register page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// Чтение данных из запроса
	var user User
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Проверка на существование пользователя
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": user.Email})
	if err != nil || count > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Сохранение пользователя в базу данных
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Registration successful!")
}

// loginHandler обработчик для логина
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("static/login.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// Чтение данных из запроса
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Поиск пользователя в базе данных
	var user User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Сравнение пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Сохранение email в сессию
	session, _ := sessionStore.Get(r, "user-session")
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	// Перенаправление на index.html
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// initAuth инициализация коллекции пользователей
func initAuth(dbClient *mongo.Client) {
	userCollection = dbClient.Database("Shop").Collection("users")
}
