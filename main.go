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
	"github.com/robfig/cron"
	"github.com/spf13/viper"

	"github.com/LiamPimlott/lunchmore/invites"
	"github.com/LiamPimlott/lunchmore/mail"
	auth "github.com/LiamPimlott/lunchmore/middleware"
	"github.com/LiamPimlott/lunchmore/organizations"
	"github.com/LiamPimlott/lunchmore/scheduling"
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
	c          *cron.Cron
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

	c = cron.New()

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

	schedulingRepository := scheduling.NewMysqlSchedulingRepository(db)
	orgsRepository := organizations.NewMysqlOrganizationsRepository(db)
	usersRepository := users.NewMysqlUsersRepository(db)
	invitesRepository := invites.NewMysqlInvitationsRepository(db)

	//////////////
	// Services //
	//////////////

	mailService := mail.NewMailService(mailConfig)
	orgsService := organizations.NewOrganizationsService(orgsRepository, mailService)
	usersService := users.NewUsersService(usersRepository, orgsService, secret)
	schedulingService := scheduling.NewSchedulingService(mailService, usersService, schedulingRepository)
	invitesService := invites.NewInviteService(invitesRepository, mailService, orgsService, usersService)

	//test mail.
	// mailService.SendText("Hello Mail.", []string{"liam.tj.pimlott@gmail.com"})

	//////////////
	// Handlers //
	//////////////

	loginUserHandler := users.NewLoginHandler(usersService)
	signupHandler := users.NewSignupHandler(usersService)

	createOrganizationHandler := organizations.NewCreateOrganizationHandler(orgsService)

	sendInviteHandler := invites.NewSendInviteHandler(invitesService)
	acceptInviteHandler := invites.NewAcceptInviteHandler(invitesService)
	getInviteHandler := invites.NewGetInviteHandler(invitesService)

	//////////////////
	// Cron Startup //
	//////////////////

	err = schedulingService.ScheduleAll(c)
	if err != nil {
		log.Fatal(err)
	}

	c.Start()
	defer c.Stop()

	////////////
	// Routes //
	////////////

	r := mux.NewRouter()

	// Users
	r.Handle("/signup", signupHandler).Methods("POST")
	r.Handle("/users/login", loginUserHandler).Methods("POST")

	// Organizations
	r.Handle("/organization", auth.Required(createOrganizationHandler, secret)).Methods("POST")

	// Invitations
	r.Handle("/invite", auth.Required(sendInviteHandler, secret)).Methods("POST")
	r.Handle("/invite", getInviteHandler).Methods("GET")
	r.Handle("/invite/accept", acceptInviteHandler).Methods("POST")

	// TODO: expose more info for invite frontend to display
	// r.Handle("/organization/invite/{code}", getInviteHandler).Methods("GET")

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
