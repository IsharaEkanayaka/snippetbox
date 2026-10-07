package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

//define an aaplication struct to hold the application-wide dependencies for the
// web application. for now we'll only include fields for the 2 custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
}

func main(){
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	//importantly we use the flag.Parse() function to parse the command-line flag.
	// this reads in the command-line flag value and assigns it to the addr 
	//variable. you need to call this before you use the addr variable
	//otherwise it will always contain the default value of ":4000", if any errors are
	//encountered during parsing the application will be terminated
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
	}


	//initialize a new http.Server struct. we set the addr and handleer fields ao
	//that the server uses network address and routes as before, and set
	// the errorlog filed so that server now uses the custom errorlog logger in
	// the event of any problem
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}