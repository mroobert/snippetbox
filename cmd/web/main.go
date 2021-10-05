package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mroobert/snippetbox/pkg/models/mysql"
)

func main() {
	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	//The parseTime=true part of the DSN is a driver-specific parameter
	//which instructs our driver to convert SQL TIME and DATE fields to Go time.Time objects.
	dsn := flag.String("dsn", "root:root@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")

	// The role of the flag.Parse() function is to parse the command-line flag.
	// This reads in the command-line flag values and assigns them to the
	// variables above. You need to call this *before* you use the variables
	// otherwise it will always contain the default value. If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Initialize a db connections pool
	db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }

    // We also defer a call to db.Close(), so that the connection pool is closed
    // before the main() function exits.
    defer db.Close()

	// Initialize a new template cache...
    templateCache, err := newTemplateCache("./ui/html/")
    if err != nil {
        errorLog.Fatal(err)
    }

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the network address from command line and the routes we defined, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger(and not the standard one) in
	// the event of any problems.
	httpServer := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = httpServer.ListenAndServe()
	errorLog.Fatal(err)
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}


// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
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
