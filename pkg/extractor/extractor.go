package extractor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kn100/eufyExtractor/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Extractor struct {
	SqlDB *sql.DB
}

type measurement struct {
	Name         string // What is this measurement
	TotalCols    int    // How many columns the display has
	ColNum       int    // Which column is this measurement on? (start at 0)
	RowNum       int    // Which row is this measurement on? (start at 1)
	YPosOverride int    // If this is a special metric (like the date), override the Ypos.
	unit         string // What to append to the end of the measurement
	whitelist    string // What characters this measurement is expected to contain
}

type ExtractorResults struct {
	Measurements []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"measurements"`
	Date int64 `json:"date"`
}

func getValFromExtractorResults(e ExtractorResults, mtype string) float64 {
	for _, elem := range e.Measurements {
		if elem.Type == mtype {
			fl, err := strconv.ParseFloat(elem.Value, 64)
			if err != nil {
				panic("was unable to parse a value as float.")
			}
			return fl
		}
	}
	return 0
}

func (e *Extractor) ProcessResultsFromExtractor(res string) error {
	exResult := ExtractorResults{}
	resBytes := []byte(res)

	err := json.Unmarshal(resBytes, &exResult)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dateStr := time.Unix(exResult.Date, 0).Format(time.RFC3339)
	if e.MeasurementExists(dateStr) {
		fmt.Println("Already seen this measurement. Ignoring.")
		return nil
	}
	sc := models.ScaleResult{

		Date:                 dateStr,
		Weight:               getValFromExtractorResults(exResult, "weight"),
		BodyFatPercentage:    getValFromExtractorResults(exResult, "body_fat_percentage"),
		Bmi:                  getValFromExtractorResults(exResult, "body_mass_index"),
		WaterPercentage:      getValFromExtractorResults(exResult, "water_percentage"),
		MuscleMassPercentage: getValFromExtractorResults(exResult, "muscle_mass_percentage"),
		BoneMassPercentage:   getValFromExtractorResults(exResult, "bone_mass_percentage"),
		BasalMetabolicRate:   getValFromExtractorResults(exResult, "basal_metabolic_rate"),
		VisceralFat:          getValFromExtractorResults(exResult, "visceral_fat"),
		LeanBodyMass:         getValFromExtractorResults(exResult, "lean_body_mass"),
		BodyFatMass:          getValFromExtractorResults(exResult, "body_fat_mass"),
		BoneMass:             getValFromExtractorResults(exResult, "bone_mass"),
		MuscleMass:           getValFromExtractorResults(exResult, "muscle_mass"),
		BodyAge:              getValFromExtractorResults(exResult, "body_age"),
		ProteinPercentage:    getValFromExtractorResults(exResult, "protein_percentage"),
	}
	// TODO: handle error
	sc.Insert(context.Background(), e.SqlDB, boil.Infer())

	return nil
}

func (e *Extractor) MeasurementExists(date string) bool {
	_, err := models.ScaleResults(qm.Where("date=?", date)).One(context.Background(), e.SqlDB)
	if err != nil {
		return false
	}
	return true
}
