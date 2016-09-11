package start

import (
	"io"
	"os"
	"strconv"
	SYS "syscall"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/rest"
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

	port, err := GetEnvInt("SERVER_PORT")
	if err != nil {
		log.Critical("Unable to parse port in .env file: ", err)
		os.Exit(1)
	}

	restService := rest.NewRestService(port)
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
