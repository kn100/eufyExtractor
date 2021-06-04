# eufyExtractor
OCR thing for extracting your data out of EufyLife digital scale apps. Very prototype.

## What it does
EufyLife is an app supplied by Eufy that connects to your smart scales. It records various measurements about your body, but annoyingly does not let you extract said information. This project is a hacky script which expects you to supply it a full screenshot (see example.jpg I have supplied), and it will do its best to extract the info using OCR. It's extra hacky, and the code isn't great, but hopefully someone else gets use out of it!

## How to use it
You'll need to run this on a system that already has Tesseract installed. On Ubuntu, this is as simple as
```bash
sudo apt-get install -yq \
    libtesseract-dev \
    libleptonica-dev
```

Then, you'll need to adjust the consts to suit your device. It's pretty simple, all you need to do is:
1. measure in pixels from the top of your screenshot to where the top of the date is, and set `DateRowXPos`
2. measure in pixels from the top of your screenshot to where the bottom of the coloured line appears in the first row of metrics, and set `FirstMetricRowXPos`
3. measure in pixels from the top of your screenshot to where the bottom of the coloured line appears in the second row of metrics, and set `SecondMetricRowXPos`
4. Measure the distance in pixels between the top of the coloured line you measured in 2 to the top of the {high, normal, low} text that appears below the metric, and set
`MetricRowYHeight`
5. If your scale displays less or more metrics, or displays them in a different order, you might need to modify the `measurements` slice.

Then, `go mod download` and `go run main.go your-screenshot.jpg with-headers`. Hopefully, you get a nice CSV.

## How I use this
I have a system set up so that every morning, I weigh myself, and then take a full screenshot of the output. My screenshot is then automatically uploaded by Nextcloud to my server. The server has a cron job set up with some bash automation that looks for new screenshots, and processes them through this, and writes it out to a CSV.

## Why?
Honestly, writing code to help with weight loss is easier than actually losing weight...  ( ͡° ͜ʖ ͡°)
