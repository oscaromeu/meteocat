package meteocat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenData[]struct {
	ID           string `json:"id"`
	CodiEstacio  string `json:"codi_estacio"`
	CodiVariable string `json:"codi_variable"`
	DataLectura  string `json:"data_lectura"`
	DataExtrem   string `json:"data_extrem,omitempty"`
	ValorLectura string `json:"valor_lectura"`
	CodiBase     string `json:"codi_base"`
}

type OpenDataMeasurements struct {
	OpenData
	*Settings
}


// NewMesurades returns a new MesuradesData pointer with the supplied parameters
func NewOpenDataMesurades() (*OpenDataMeasurements, error) {
	c := &OpenDataMeasurements{
		Settings: NewSettings(),
	}


	return c, nil
}

func (s *OpenDataMeasurements) OpenDataMeasurementAllByStation(p *Parameters) error {


	req, err := http.NewRequest("GET", fmt.Sprintf(openDataURL, fmt.Sprintf("?data_lectura=%s-%s-%sT06:00:00", p.Any, p.Mes, p.Dia)), nil)
	if err != nil {
		return err
	}

	//req.Header.Add("X-Api-Key", m.Key)

	resp, err := s.client.Do(req)


	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&s.OpenData); err != nil {
		return err
	}

	return nil
}


