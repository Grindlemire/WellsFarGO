package start

import (
	"io"
	"os"
	"strconv"
	SYS "syscall"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/rest"
	"github.com/grindlemire/WellsFarGO/unifier"
	"github.com/joho/godotenv"
	DEATH "github.com/vrecan/death"
)

// Run starts the webserver
func Run() {
	var goRoutines []io.Closer
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		log.Critical("Error loading .env file")
		os.Exit(1)
	}

	dbFile := os.Getenv("DB_FILE")
	csvFile := os.Getenv("CSV_FILE")
	formatType := os.Getenv("FORMAT_TYPE")

	unifier, err := unifier.NewUnifier(dbFile, csvFile, formatType)
	if err != nil {
		log.Critical("Error initializing the unifier")
		os.Exit(1)
	}

	err = unifier.AddNewData()
	if err != nil {
		log.Critical("Could not add new data: ", err)
		os.Exit(1)
	}

	// sTime, _ := time.Parse("01/02/2006", "06/13/2016")
	// eTime, _ := time.Parse("01/02/2006", "06/14/2016")
	// results, _ := unifier.QueryDateRange(sTime, eTime)
	// fmt.Printf("Date Range Query: %#v\n\n", results)
	//
	// dTime, _ := time.Parse("01/02/2006", "06/13/2016")
	// results, _ = unifier.QueryDay(dTime)
	// fmt.Printf("Day Query: %#v\n\n", results)
	//
	// results, _ = unifier.QueryAmountRange(0.0, 10.00)
	// fmt.Printf("Amount Range Query: %#v\n\n", results)
	//
	// results, _ = unifier.QueryAmount(1.06)
	// fmt.Printf("Amount Query: %#v\n\n", results)
	//
	// results, _ = unifier.QueryLocation("TACO BELL 302600302646 BROOMFIELD CO")
	// fmt.Printf("Amount Query: %#v\n\n", results)

	// WEBSERVER
	port, err := GetEnvInt("SERVER_PORT")
	if err != nil {
		log.Critical("Unable to parse port in .env file: ", err)
		os.Exit(1)
	}

	restService := rest.NewRestService(port, unifier)
	goRoutines = append(goRoutines, restService)
	restService.Start()

	death.WaitForDeath(goRoutines...)

}

// GetEnvInt gets an environment variable and returns it as an int
func GetEnvInt(key string) (val int, err error) {
	strVal := os.Getenv(key)
	val, err = strconv.Atoi(strVal)
	return val, err
}
