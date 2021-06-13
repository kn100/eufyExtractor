FROM golang:1.16.4-buster as builder

# Install some dependencies needed to build the project
# not sure why I need libgl-dev but whatever.
RUN apt-get update && apt-get install -y bash ca-certificates git gcc g++ libc-dev tesseract-ocr libgl-dev python3 python3-pip

RUN mkdir /app
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
ENV GO111MODULE=on
RUN go mod download

COPY . /app
RUN go build -o main .

# Build python
WORKDIR /app/python-extractor
# Not sure why I need to install Pillow outside, but I do.
RUN pip3 install Pillow
RUN pip3 install -r requirements.txt
RUN pip3 install pyinstaller
RUN pyinstaller --onefile main.py
RUN cp dist/main ../magic-extractor

WORKDIR /app
CMD ["/app/main"]
