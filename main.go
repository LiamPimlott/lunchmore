package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/LiamPimlott/lunchmore/mail"
	auth "github.com/LiamPimlott/lunchmore/middleware"
	"github.com/LiamPimlott/lunchmore/organizations"
	"github.com/LiamPimlott/lunchmore/users"
)

// DatabaseConfig describes configuration for connecting to a database
type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

var (
	dbConfig   DatabaseConfig
	mailConfig mail.ClientConfig
	secret     string
	serverPort string
)

func init() {
	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	serverPort = viper.GetString(`server.port`)
	dbConfig = DatabaseConfig{
		Host: viper.GetString(`database.host`),
		Port: viper.GetString(`database.port`),
		User: viper.GetString(`database.user`),
		Pass: viper.GetString(`database.pass`),
		Name: viper.GetString(`database.name`),
	}
	mailConfig = mail.ClientConfig{
		Username: viper.GetString(`mail.username`),
		Password: viper.GetString(`mail.password`),
		Host:     viper.GetString(`mail.host`),
		Port:     viper.GetInt(`mail.port`),
	}
	secret = viper.GetString(`jwt.secret`)
}

func main() {
	////////
	// DB //
	////////

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected")

	///////////
	// Repos //
	///////////

	orgsRepository := organizations.NewMysqlOrganizationsRepository(db)
	usersRepository := users.NewMysqlUsersRepository(db)

	//////////////
	// Services //
	//////////////

	mailService := mail.NewMailService(mailConfig)
	orgsService := organizations.NewOrganizationsService(orgsRepository, secret)
	usersService := users.NewUsersService(usersRepository, secret)

	//test mail.
	mailService.SendText("Hello Mail.", []string{"liam.tj.pimlott@gmail.com"})

	//////////////
	// Handlers //
	//////////////

	loginUserHandler := users.NewLoginHandler(usersService)
	createUserHandler := users.NewCreateUserHandler(usersService)
	createOrganizationHandler := organizations.NewCreateOrganizationHandler(orgsService)
	// getUserByIDHandler := users.NewGetUserByIDHandler(usersService)

	////////////
	// Routes //
	////////////

	r := mux.NewRouter()

	// Users

	r.Handle("/users", createUserHandler).Methods("POST")
	r.Handle("/users/login", loginUserHandler).Methods("POST")
	// r.Handle("/users/{id}", auth.Required(getUserByIDHandler, secret)).Methods("GET")

	// Organizations
	r.Handle("/organizations", auth.Required(createOrganizationHandler, secret)).Methods("POST")

	// TODO /organizations POST
	// TODO /organizations/{id}/user POST

	////////////
	// STATIC //
	////////////

	// serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))

	// root route will serve the built react app
	r.Handle("/", http.FileServer(http.Dir("./frontend/build")))

	// start server
	log.Printf("listening on port %s\n", serverPort)
	http.ListenAndServe(serverPort, handlers.LoggingHandler(os.Stdout, r))
}
