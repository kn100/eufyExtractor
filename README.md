# eufyExtractor
OCR thing for extracting your data out of EufyLife digital scale apps. Very
prototype.

## What it does
EufyLife is an app supplied by Eufy that connects to your smart scales. It
records various measurements about your body, but annoyingly does not let you
extract said information. This project is a hacky thing which expects you to
supply it a full screenshot (see example.jpg I have supplied), and it will do
its best to extract the info using OCR. It's extra hacky, and the code isn't
great, but hopefully someone else gets use out of it!

See the related blog post: https://kn100.me/taking-back-data-from-eufy/

## How to use it

Unless you specifically have a Huawei P30, you'll probably need to tweak the
constants in `python-extractor/main.py`. See `exampleImage.png` for how to get
these numbers. You'll also want to tweak the Filename regex in the Dockerfile.

I suggest using the supplied Dockerfile + Docker Compose file. You'll want to
change the volume mounts in the docker-compose file to match your own import
directory and where you'd like to store the Sqlite database.

```
docker-compose build
docker-compose up -d
docker-compose logs --follow eufyextractor_backend_1
```

Once your service is up, you can query port 52525 and will get data as JSON.

Good luck if you attempt to use this outside of a container :P
## How I use this
I have a system set up so that every morning, I weigh myself, and then take a
full screenshot of the output. This program then detects that new screenshot,
processes it, and I've got a frontend set up at kn100.me/weight-loss/ which
pulls the data from the endpoint in the service.

## Why?
Honestly, writing code to help with weight loss is easier than actually losing
weight...  ( ͡° ͜ʖ ͡°)

## TODO:
* Expose the metrics as environment variables, rather than requiring the user to
  modify code
* Improve the logging across the board
* Figure out how to make the Docker build and run containers separate
* Add ways to filter the data the service returns (for example, by time range)
* Write a frontend I am happy to share
* Potentially rewrite the entire lot in Python
* Tidy up the hideous date parsing logic
