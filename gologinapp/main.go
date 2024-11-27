package main  

import 
(  
	"crypto/rand"  
	"database/sql"  
	"encoding/base64"  
	"encoding/json"  
	"fmt"  
	"log"  
	"net/http"  
	"os"  
	"time"  

	"github.com/dgrijalva/jwt-go"  
	"github.com/lib/pq"  
	"golang.org/x/crypto/bcrypt"  
)  

var db *sql.DB  
var jwtSecretKey string  
 
type User struct 
{  
	ID       int    `json:"id"`  
	Username string `json:"username"`  
	Password string `json:"-"`  
}  
 
type Claims struct 
{  
	Username string `json:"username"`  
	jwt.StandardClaims  
}  

func init() 
{  
	var err error  
	connStr := "user=postgres dbname=login password=root sslmode=disable" 
	db, err = sql.Open("postgres", connStr)  
	if err != nil 
	{  
		log.Fatal(err)  
	}  
 
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")  
	if jwtSecretKey == "" 
	{  
		jwtSecretKey = generateSecretKey()  
		fmt.Println("Generated new secret key:", jwtSecretKey)  
		  
	}  
}  

func generateSecretKey() string 
{  
	key := make([]byte, 32)   
	if _, err := rand.Read(key); err != nil 
	{  
		log.Fatal("Error generating secret key:", err)  
	}  
	return base64.StdEncoding.EncodeToString(key) 
}  

func main() 
{  
	  
	http.HandleFunc("/register", register)  
	http.HandleFunc("/login", login)  
	http.HandleFunc("/welcome", welcome)  

	log.Println("Server started at :3000")  
	log.Fatal(http.ListenAndServe(":3000", nil))  
}  



func register(w http.ResponseWriter, r *http.Request) 
{  
	if r.Method != http.MethodPost 
	{  
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)  
		return  
	}  

	var user User  
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil 
	{  
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return  
	}  

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)  
	if err != nil 
	{  
		http.Error(w, err.Error(), http.StatusInternalServerError)  
		return  
	}  

	user.Password = string(hashedPassword)  

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)  
	if err != nil 
	{  
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" 
		{  
			http.Error(w, "Username already exists", http.StatusConflict)  
		} else 
		{  
			http.Error(w, err.Error(), http.StatusInternalServerError)  
		}  
		return  
	}  
	w.WriteHeader(http.StatusCreated)  
}  

func login(w http.ResponseWriter, r *http.Request) 
{  
	if r.Method != http.MethodPost 
	{  
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)  
		return  
	}  

	var user User  
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil 
	{  
		http.Error(w, err.Error(), http.StatusBadRequest)  
		return  
	}  

	var storedUser User  
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", user.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)  
	if err != nil 
	{  
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)  
		return  
	}  

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))  
	if err != nil 
	{  
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)  
		return  
	}  

	claims := &Claims
	{  
		Username: storedUser.Username,  
		StandardClaims: jwt.StandardClaims
		{  
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),  
		},  
	}  

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  
	tokenString, err := token.SignedString([]byte(jwtSecretKey)) 
	if err != nil 
	{  
		http.Error(w, err.Error(), http.StatusInternalServerError)  
		return  
	}  

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})  
}  

func welcome(w http.ResponseWriter, r *http.Request) 
{  
	tokenString := r.Header.Get("Authorization")  
	if tokenString == "" 
	{  
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)  
		return  
	}  

	claims := &Claims{}  
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) 
	{  
		return []byte(jwtSecretKey), nil 
	})  

	if err != nil || !token.Valid 
	{  
		http.Error(w, "Invalid token", http.StatusUnauthorized)  
		return  
	}  

	fmt.Fprintf(w, "Welcome %s!", claims.Username)  
}