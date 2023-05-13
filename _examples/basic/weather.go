package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oscaromeu/meteocat"
)

func main() {

	// execute export METEOCAT_API_KEY=<API_KEY_VALUE> on a shell first
	d, err := meteocat.NewMesurades(os.Getenv("METEOCAT_API_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	if meteocat.CheckAPIKeyExists(d.Key) == false {
		fmt.Println("ApiKey is not set. ")
	}

	data := meteocat.Data{
		Any: "2023",
		Mes: "01",
		Dia: "06",
	}
	params, _ := meteocat.NewParameters(
		//meteocat.OptionCodiEstacio("D5"),
		meteocat.OptionCodiVariable("32"),
		meteocat.OptionData(data),
	)

	// Call MeasurementByDay Method
	err = d.MeasurementByDay(params)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(d.Measurements)
}
