package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kn100/eufyExtractor/modext"
	"github.com/kn100/eufyExtractor/pkg/extractor"
	_ "github.com/mattn/go-sqlite3"
	"github.com/radovskyb/watcher"
)

var db *sql.DB
var dbPath string = "./database/eufyExtractor-data.db"

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("Starting eufyExtractor")
	w := watcher.New()
	w.FilterOps(watcher.Create)

	r := regexp.MustCompile(EnvString("FILENAME_REGEX", ".+"))
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	var err error
	createDbIfNotExist()
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err) //TODO
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.WithField("image", event.Path).Infoln("Processing image")

				err = handleImage(db, event.Path)
				if err != nil {
					log.Errorln(err)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(EnvString("IMPORT_DIR", "import")); err != nil {
		log.Fatalln(err)
	}
	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
	http.HandleFunc("/", measurements)
	http.ListenAndServe(":"+EnvString("HTTP_PORT", "52525"), nil)

}

func measurements(w http.ResponseWriter, req *http.Request) {
	fromString := req.URL.Query().Get("from")
	toString := req.URL.Query().Get("to")
	from := time.Time{}
	if fromString != "" {
		var err error
		from, err = time.Parse(time.RFC3339, fromString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}
	}
	to := time.Now()
	if toString != "" {
		var err error
		to, err = time.Parse(time.RFC3339, toString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}
	}
	json := modext.GetPercsAsJSON(db, from, to)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, json)
}

func handleImage(db *sql.DB, path string) error {
	out, err := exec.Command("./magic-extractor", path).Output()
	if err != nil {
		log.Errorln("something went wrong running command " + err.Error())
		return err
	}
	outString := string(out)
	if strings.HasPrefix(outString, "error") {
		log.Errorln("Error from extractor: " + outString)
	}

	e := extractor.Extractor{
		SqlDB: db,
	}
	err = e.ProcessResultsFromExtractor(outString)
	if err != nil {
		log.Errorln("error processing results from image.")
		return nil
	}
	return nil
}

func EnvString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func createDbIfNotExist() error {

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		if err != nil {
			log.Fatal("Could not write file for some reason", err)
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE "scale_results" (
		id integer primary key autoincrement,
		date varchar(255) NOT NULL,
		weight REAL NOT NULL,
		bmi REAL NOT NULL,
		body_fat_percentage REAL NOT NULL,
		water_percentage REAL NOT NULL,
		muscle_mass_percentage REAL NOT NULL,
		bone_mass_percentage REAL NOT NULL,
		basal_metabolic_rate REAL NOT NULL,
		visceral_fat REAL NOT NULL,
		lean_body_mass REAL NOT NULL,
		body_fat_mass REAL NOT NULL,
		bone_mass REAL NOT NULL,
		muscle_mass REAL NOT NULL,
		body_age REAL NOT NULL,
		protein_percentage REAL NOT NULL
	)`)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
