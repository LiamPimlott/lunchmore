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
	"github.com/gorilla/sessions"
	"github.com/robfig/cron"
	"github.com/spf13/viper"

	"github.com/LiamPimlott/lunchmore/invites"
	libsessions "github.com/LiamPimlott/lunchmore/lib/sessions"
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
	c                *cron.Cron
	dbConfig         DatabaseConfig
	mailConfig       mail.ClientConfig
	secret           string
	serverPort       string
	frontendHostName string
)

func init() {
	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	c = cron.New()

	serverPort = viper.GetString(`server.port`)

	frontendHostName = viper.GetString(`frontend.hostname`)

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

	//////////////
	// SESSIONS //
	//////////////
	cookieStore := sessions.NewCookieStore([]byte("34-or-64-bits-recommended"), []byte("16-24-&-34-bytes"))

	cookieStore.Options = &sessions.Options{
		MaxAge: 7200,
		// Secure:   true,
		HttpOnly: true,
		Path:     "/users/refresh",
		// SameSite: http.SameSiteStrictMode,
	}

	sessionStorer := libsessions.NewSessionStorer(cookieStore, "lunchmore-session")

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
	schedulingService := scheduling.NewSchedulingService(c, mailService, usersService, orgsService, schedulingRepository)
	invitesService := invites.NewInviteService(invitesRepository, mailService, orgsService, usersService)

	//test mail.
	// mailService.SendText("Hello Mail.", []string{"liam.tj.pimlott@gmail.com"})

	//////////////
	// Handlers //
	//////////////

	// Users
	signupHandler := users.NewSignupHandler(usersService, sessionStorer)
	loginUserHandler := users.NewLoginHandler(usersService, sessionStorer)
	logoutUserHandler := users.NewLogoutHandler(sessionStorer)
	refreshHandler := users.NewRefreshHandler(usersService, sessionStorer)

	// Orgs
	createOrganizationHandler := organizations.NewCreateOrganizationHandler(orgsService)

	// Invites
	sendInviteHandler := invites.NewSendInviteHandler(invitesService)
	acceptInviteHandler := invites.NewAcceptInviteHandler(invitesService)
	getInviteHandler := invites.NewGetInviteHandler(invitesService)

	// Schedules
	createScheduleHandler := scheduling.NewCreateScheduleHandler(schedulingService)
	getOrgSchedulesHandler := scheduling.NewGetOrgSchedulesHandler(schedulingService)
	joinScheduleHandler := scheduling.NewJoinScheduleHandler(schedulingService)

	//////////////////
	// Cron Startup //
	//////////////////

	err = schedulingService.ScheduleAll()
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
	r.Handle("/users/logout", logoutUserHandler).Methods("GET")
	r.Handle("/users/refresh", refreshHandler).Methods("GET")

	// Orgs
	r.Handle("/organization", auth.Required(createOrganizationHandler, secret)).Methods("POST")

	// Invitations
	r.Handle("/invite", auth.Required(sendInviteHandler, secret)).Methods("POST")
	r.Handle("/invite", getInviteHandler).Methods("GET")
	r.Handle("/invite/accept", acceptInviteHandler).Methods("POST")

	// Schedules
	r.Handle("/schedules", auth.Required(createScheduleHandler, secret)).Methods("POST")
	r.Handle("/schedules/organization/{id}", auth.Required(getOrgSchedulesHandler, secret)).Methods("GET")
	r.Handle("/schedules/{id}/join", auth.Required(joinScheduleHandler, secret)).Methods("GET")

	// TODO: expose more info for invite frontend to display
	// r.Handle("/organization/invite/{code}", getInviteHandler).Methods("GET")

	////////////
	// STATIC //
	////////////

	// serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))

	// root route will serve the built react app
	r.Handle("/", http.FileServer(http.Dir("./frontend/build")))

	//////////
	// CORS //
	//////////

	corsConfig := []handlers.CORSOption{
		handlers.AllowedOrigins([]string{frontendHostName}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	}

	////////////////
	// MIDDLEWARE //
	////////////////

	// chain logging middleware
	h := handlers.LoggingHandler(os.Stdout, r)

	// chain CORS middleware
	h = handlers.CORS(corsConfig...)(h)

	// start server
	log.Printf("listening on port %s\n", serverPort)
	http.ListenAndServe(serverPort, h)
}
