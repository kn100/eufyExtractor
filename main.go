package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/oliamb/cutter"
	"github.com/otiai10/gosseract/v2"
)

// The following values are based on a Huawei P30. You may need to adjust these
// for your device assuming the UI is similar.
const DateRowXPos = 250
const FirstMetricRowXPos = 663
const SecondMetricRowXPos = 1253
const MetricRowYHeight = 125

// These control the characters we are looking for. You might need to adjust the
// date if you aren't using English.
const whitelistMeasurement = ".1234567890"
const whitelistDate = " .:0123456789aAbcDeFgJlMnNoOpPrStuvy"

type measurement struct {
	Name         string // What is this measurement
	TotalCols    int    // How many columns the display has
	ColNum       int    // Which column is this measurement on? (start at 0)
	RowNum       int    // Which row is this measurement on? (start at 1)
	YPosOverride int    // If this is a special metric (like the date), override the Ypos.
	unit         string // What to append to the end of the measurement
	whitelist    string // What characters this measurement is expected to contain

}

type completedMeasurement struct {
	Name  string
	Value string
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("You need to specify an image.")
		os.Exit(1)
	}

	measurements := []measurement{
		{Name: "Date", TotalCols: 1, ColNum: 0, YPosOverride: DateRowXPos, whitelist: whitelistDate},
		{Name: "Weight", TotalCols: 2, ColNum: 0, RowNum: 0, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "BMI", TotalCols: 2, ColNum: 1, RowNum: 0, whitelist: whitelistMeasurement, unit: "%"},
		{Name: "Body Fat %", TotalCols: 2, ColNum: 0, RowNum: 1, whitelist: whitelistMeasurement, unit: "%"},
		{Name: "Water", TotalCols: 2, ColNum: 1, RowNum: 1, whitelist: whitelistMeasurement},
		{Name: "Muscle Mass %", TotalCols: 2, ColNum: 0, RowNum: 2, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "Bone Mass %", TotalCols: 2, ColNum: 1, RowNum: 2, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "BMR", TotalCols: 2, ColNum: 0, RowNum: 3, whitelist: whitelistMeasurement},
		{Name: "Visceral Fat", TotalCols: 2, ColNum: 1, RowNum: 3, whitelist: whitelistMeasurement},
		{Name: "Lean Body Mass", TotalCols: 2, ColNum: 0, RowNum: 4, whitelist: whitelistMeasurement, unit: "%"},
		{Name: "Body Fat Mass", TotalCols: 2, ColNum: 1, RowNum: 4, whitelist: whitelistMeasurement, unit: "%"},
		{Name: "Bone Mass", TotalCols: 2, ColNum: 0, RowNum: 5, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "Muscle Mass", TotalCols: 2, ColNum: 1, RowNum: 5, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "Body Age", TotalCols: 2, ColNum: 0, RowNum: 6, whitelist: whitelistMeasurement, unit: "Kg"},
		{Name: "Protein", TotalCols: 2, ColNum: 1, RowNum: 6, whitelist: whitelistMeasurement, unit: "%"},
	}

	client := gosseract.NewClient()
	defer client.Close()

	img, err := getImage(os.Args[1])
	if err != nil {
		panic(err)
	}

	imgWidth, _, _, err := getImageData(os.Args[1])
	if err != nil {
		panic(err)
	}
	// Compute how far we need to travel to get to next row
	DistBetweenMetricRows := SecondMetricRowXPos - FirstMetricRowXPos

	// Where we store our final measurements
	var completedMeasurements []completedMeasurement

	for _, element := range measurements {
		xpos := getBoxX(imgWidth, element)
		ypos := getBoxY(element, DistBetweenMetricRows, FirstMetricRowXPos)
		colWidth := int(math.Floor(float64(imgWidth) / float64(element.TotalCols)))

		//fmt.Printf("Cropping %s at x: %d, y:%d, w:%d h:%d\n", element.Name, xpos, ypos, colWidth, MetricRowYHeight)

		// First, crop the image around the metric in question.
		outImage, err := cutter.Crop(img, cutter.Config{
			Width:  colWidth,
			Height: MetricRowYHeight,
			Anchor: image.Point{xpos, ypos},
		})

		// Next, apply some post processing magic to improve OCR accuracy
		isDate := element.Name == "Date"
		outImageNRGB, err := processImage(outImage, !isDate)
		if err != nil {
			panic(err)
		}

		// Next, convert it to a PNG for Tesseract
		buf := new(bytes.Buffer)
		err = png.Encode(buf, &outImageNRGB)
		if err != nil {
			panic(err)
		}

		// Below line writes the processed images out to /tmp/ so you can review them.
		// ioutil.WriteFile("/tmp/"+element.Name+".png", buf.Bytes(), 0644)

		// Do the OCR
		client.SetWhitelist(element.whitelist)
		client.SetImageFromBytes(buf.Bytes())
		text, err := client.Text()
		if err != nil {
			panic(err)
		}

		if isDate {
			parsedDate, err := tryParseDate(text)
			if err != nil {
				fmt.Println(err)
			}
			text = parsedDate
		}

		// Post process the string somewhat, and append it to completed measurements
		text = processString(text)
		text = text + element.unit
		cm := completedMeasurement{Name: element.Name, Value: text}
		completedMeasurements = append(completedMeasurements, cm)
	}

	// Write the output as a CSV.
	var headers string
	var values string
	for i, cma := range completedMeasurements {
		if i+1 == len(completedMeasurements) {
			headers = headers + cma.Name
			values = values + cma.Value
		} else {
			headers = headers + cma.Name + ", "
			values = values + cma.Value + ", "
		}

	}
	if len(os.Args) == 3 && os.Args[2] == "with-headers" {
		fmt.Println(headers)
	}
	fmt.Println(values)
}

func tryParseDate(date string) (string, error) {
	formats := []string{
		"Jan.02 2006 03:04:PM",
		"Jan.02 2006 3:04:PM",
		"Jan.2 2006 03:04:PM",
		"Jan.2 200603:04:PM",
	}
	for _, format := range formats {
		t, err := time.Parse(format, date)
		if err == nil {
			return t.Format("02/01/2006 15:04"), nil
		}
	}
	return "", fmt.Errorf("%s Date wasn't parseable", date)
}

func getImage(imagePath string) (image.Image, error) {
	dat, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(dat)

	return img, err
}

func getImageData(imagePath string) (int, int, string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, "", fmt.Errorf("%v\n", err)
	}

	image, format, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, "", fmt.Errorf("%s: %v\n", imagePath, err)
	}

	return image.Width, image.Height, format, nil
}

func processImage(img image.Image, shouldProcess bool) (image.NRGBA, error) {
	// Only here to convert image.Image to image.NRGBA.
	image := imaging.AdjustBrightness(img, 0)

	if shouldProcess {
		image = imaging.Blur(image, 3.5)
		image = imaging.AdjustBrightness(image, 9)
		image = imaging.AdjustContrast(image, 100)
		image = imaging.Sharpen(image, 250)
	}

	return *image, nil
}

func processString(in string) string {
	out := strings.Trim(in, ".")
	return strings.ReplaceAll(out, " ", "")
}

func getBoxX(screen_width int, measurement measurement) int {
	colWidth := int(math.Floor(float64(screen_width) / float64(measurement.TotalCols)))
	return colWidth * measurement.ColNum
}

func getBoxY(measurement measurement, dist_between_rows, topBoxFirstRowPos int) int {
	if measurement.YPosOverride != 0 {
		return measurement.YPosOverride
	}
	return topBoxFirstRowPos + (measurement.RowNum * dist_between_rows)
}
