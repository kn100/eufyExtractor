package modext

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/kn100/eufyExtractor/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type datasets struct {
	Data []dataset `json:"datasets"`
}

type dataset struct {
	Label string
	Data  []WeightReturn `json:"data"`
}

type WeightReturn struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

func GetPercsAsJSON(db *sql.DB, from, to time.Time) string {
	// time filtering is broken, so I've disabled it for now.
	scaleResults, err := models.ScaleResults(qm.OrderBy("date")).
		All(context.Background(), db)
	if err != nil {
		panic(err)
	}
	weight := []WeightReturn{}
	muscleMassPerc := []WeightReturn{}
	waterPerc := []WeightReturn{}
	bodyFatPerc := []WeightReturn{}
	proteinPerc := []WeightReturn{}
	boneMassPerc := []WeightReturn{}
	basalMetabolicRate := []WeightReturn{}
	bodyAge := []WeightReturn{}
	bmi := []WeightReturn{}
	bodyFatMass := []WeightReturn{}
	boneMass := []WeightReturn{}
	leanBodyMass := []WeightReturn{}
	muscleMass := []WeightReturn{}
	visceralFat := []WeightReturn{}
	for _, el := range scaleResults {
		measurementTime, err := time.Parse(time.RFC3339Nano, el.Date)
		if err != nil {
			panic(err)
		}
		measurementTimeUnix := measurementTime.Unix() * 1000
		weight = append(weight, WeightReturn{X: measurementTimeUnix, Y: el.Weight})
		muscleMassPerc = append(muscleMassPerc, WeightReturn{X: measurementTimeUnix, Y: el.MuscleMassPercentage})
		waterPerc = append(waterPerc, WeightReturn{X: measurementTimeUnix, Y: el.WaterPercentage})
		bodyFatPerc = append(bodyFatPerc, WeightReturn{X: measurementTimeUnix, Y: el.BodyFatPercentage})
		proteinPerc = append(proteinPerc, WeightReturn{X: measurementTimeUnix, Y: el.ProteinPercentage})
		boneMassPerc = append(boneMassPerc, WeightReturn{X: measurementTimeUnix, Y: el.BoneMassPercentage})
		basalMetabolicRate = append(basalMetabolicRate, WeightReturn{X: measurementTimeUnix, Y: el.BasalMetabolicRate})
		bodyAge = append(bodyAge, WeightReturn{X: measurementTimeUnix, Y: el.BodyAge})
		bmi = append(bmi, WeightReturn{X: measurementTimeUnix, Y: el.Bmi})
		bodyFatMass = append(bodyFatMass, WeightReturn{X: measurementTimeUnix, Y: el.BodyFatMass})
		boneMass = append(boneMass, WeightReturn{X: measurementTimeUnix, Y: el.BoneMass})
		leanBodyMass = append(leanBodyMass, WeightReturn{X: measurementTimeUnix, Y: el.LeanBodyMass})
		muscleMass = append(muscleMass, WeightReturn{X: measurementTimeUnix, Y: el.MuscleMass})
		visceralFat = append(visceralFat, WeightReturn{X: measurementTimeUnix, Y: el.VisceralFat})

	}
	datasets := datasets{}
	datasets.Data = append(datasets.Data,
		dataset{Label: "weight", Data: weight},
		dataset{Label: "muscle_mass_percentage", Data: muscleMassPerc},
		dataset{Label: "water_percentage", Data: waterPerc},
		dataset{Label: "body_fat_percentage", Data: bodyFatPerc},
		dataset{Label: "protein_percentage", Data: proteinPerc},
		dataset{Label: "bone_mass_percentage", Data: boneMassPerc},
		dataset{Label: "basal_metabolic_rate", Data: basalMetabolicRate},
		dataset{Label: "body_age", Data: bodyAge},
		dataset{Label: "bmi", Data: bmi},
		dataset{Label: "body_fat_mass", Data: bodyFatMass},
		dataset{Label: "bone_mass", Data: boneMass},
		dataset{Label: "lean_body_mass", Data: leanBodyMass},
		dataset{Label: "muscle_mass", Data: muscleMass},
		dataset{Label: "visceral_fat", Data: visceralFat},
	)

	json, err := json.Marshal(datasets)
	if err != nil {
		panic(err)
	}
	return string(json)
}
