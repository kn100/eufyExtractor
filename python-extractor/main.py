import cv2
import measurements
import pytesseract
import numpy as np
import sys
import time
import logging
from PIL import Image
from datetime import datetime
DATE_ROW_Y_POS = 250
FIRST_METRIC_ROW_X_POS = 663
SECOND_METRIC_ROW_X_POS = 1253
METRIC_ROW_Y_HEIGHT = 125

SCREENSHOT_WIDTH = 1080
SCREENSHOT_HEIGHT = 5254
DIST_HEIGHT_BETWEEN_ROWS = SECOND_METRIC_ROW_X_POS - FIRST_METRIC_ROW_X_POS


def process_metric_unit_for_tesseract(image):
    original = image.copy()
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    # read image, copy, convert to greyscale, threshold edges to pick out the chars
    thresh = cv2.threshold(
        gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]

    # find bounds of each char and calculate max area and height, this does not maintain order
    chars = []
    cnts = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    cnts = cnts[0] if len(cnts) == 2 else cnts[1]
    max_area = 0
    max_height = 0
    for c in cnts:
        x, y, w, h = cv2.boundingRect(c)
        area = w * h
        max_area = max(max_area, area)
        max_height = max(max_height, h)
        chars.append((original[y:y+h, x:x+w], x, area))

    # prune all small characters except periods
    prune = []
    total_width = 0
    for c in chars:
        if not (c[2] <= max_area * 0.5 and c[2] >= max_area * 0.07):
            prune.append(c)
            total_width += len(c[0][0])

    # sort by x coordinate to put back in order
    sort = sorted(prune, key=lambda t: t[1])

    # stitch back into an image horizontally
    max_height = max_height
    char_spacing = 30
    stitch = np.zeros((max_height + 100, total_width + 130, 3), np.uint8)
    stitch[:] = (255, 255, 255)
    x = 0
    for c in sort:
        w = len(c[0][0])
        h = len(c[0])
        stitch[(max_height-h)+10:max_height+10, x+10:(x+w)+10, :3] = c[0]
        x += w+char_spacing

    stitch = cv2.cvtColor(stitch, cv2.COLOR_BGR2GRAY)
    (thresh, im_bw) = cv2.threshold(stitch, 128,
                                    255, cv2.THRESH_BINARY | cv2.THRESH_OTSU)

    kernel = np.ones((3, 3), np.uint8)
    im_bw = cv2.dilate(im_bw, kernel, iterations=1)

    return im_bw


def process_date_for_tesseract(image):
    # Thanks to the #tesseract Libera IRC channel, which has absolutely nothing to do with OpenCV or Tesseract.
    original = image.copy()
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    # read image, copy, convert to greyscale, threshold edges to pick out the chars
    thresh = cv2.threshold(
        gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]

    # find bounds of each char and calculate max area and height, this does not maintain order
    chars = []
    cnts = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    cnts = cnts[0] if len(cnts) == 2 else cnts[1]
    max_area = 0
    max_height = 0
    for c in cnts:
        # print(len(cnts))
        x, y, w, h = cv2.boundingRect(c)
        area = w * h
        max_area = max(max_area, area)
        max_height = max(max_height, h)
        chars.append((original[y:y+h, x:x+w], x, area))

    # prune all small characters except periods
    prune = []
    total_width = 0
    for c in chars:
        prune.append(c)
        total_width += len(c[0][0])

    # sort by x coordinate to put back in order
    sort = sorted(prune, key=lambda t: t[1])

    # stitch back into an image horizontally
    max_height = max_height
    char_spacing = 6
    stitch = np.zeros((max_height + 200, total_width + 200, 3), np.uint8)
    stitch[:] = (255, 255, 255)
    x = 0
    for c in sort:
        w = len(c[0][0])
        h = len(c[0])
        stitch[(max_height-h)+15:max_height+15, x+15:(x+w)+15, :3] = c[0]
        x += w+char_spacing

    stitch = cv2.cvtColor(stitch, cv2.COLOR_BGR2GRAY)
    (thresh, im_bw) = cv2.threshold(stitch, 128,
                                    255, cv2.THRESH_BINARY | cv2.THRESH_OTSU)
    return im_bw


def process_date_output(date_text, logger):
    # process_date_output processes the date output...sorry
    parsed_date_str = ""
    currpos = 0
    date_text = date_text.strip().lower()
    month_short = ['jan', 'feb', 'mar', 'apr', 'may',
                   'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']
    valid_years = [2021, 2022, 2023]
    if not (date_text[0:3] in month_short):
        logger.error("error: date was not parseable")
        exit()
    parsed_date_str = parsed_date_str + date_text[0:3]
    currpos = currpos + 3
    if date_text[3] != ".":
        logger.error("error: date was not parseable")
        exit()
    parsed_date_str = parsed_date_str + "-"
    currpos = currpos + 1
    valid_days = list(range(1, 32))

    if int(date_text[4:6]) in valid_days:
        # Not guaranteed, but probability increases this is truly the day
        if not int(date_text[5:9]) in valid_years and int(date_text[6:10]) in valid_years:
            # Pretty sure that day has two digits
            parsed_date_str = parsed_date_str + \
                date_text[4:6] + "-" + date_text[6:10]
            currpos = currpos + 6
        else:
            parsed_date_str = parsed_date_str + \
                "0" + date_text[4:5] + "-" + date_text[5:9]
            currpos = currpos + 5
    elif int(date_text[4:5]) in valid_days and int(date_text[5:9]) in valid_years:
        parsed_date_str = parsed_date_str + \
            "0" + date_text[4:5] + "-" + date_text[5:9]
        currpos = currpos + 5

    # 2 digit hour - maybe validate the hour is valid
    if date_text[currpos+2:currpos+4] == "..":
        parsed_date_str += " " + date_text[currpos:currpos+2]
        currpos = currpos + 2
    else:
        parsed_date_str += " 0" + date_text[currpos:currpos+1]
        currpos = currpos + 1

    parsed_date_str += ":" + date_text[currpos+2:currpos+4]
    currpos = currpos + 4
    parsed_date_str += "" + date_text[currpos+2:currpos+4]

    try:
        dparsed = datetime.strptime(
            parsed_date_str, '%b-%d-%Y %I:%M%p')
    except ValueError:
        logger.error("error: could not parse date")
        exit()

    return time.mktime(dparsed.timetuple())


def process_metric_output(metric_text, reasonable_range, logger):
    if len(reasonable_range) == 0:
        return metric_text
    try:
        float(metric_text)
    except ValueError:
        logger.info("Number was not a float")
        exit()
    if float(metric_text) < float(reasonable_range[0]) or float(metric_text) > float(reasonable_range[1]):
        div10 = (float(metric_text) / 10)
        if float(div10) < reasonable_range[0] or float(div10) > reasonable_range[1]:

            logger.error("error: metric outside of reasonable range")
            exit()
        else:
            return str(div10)
    else:
        return metric_text


def process_image(image, logger):
    m = measurements.Measurements(
        SCREENSHOT_WIDTH, FIRST_METRIC_ROW_X_POS, DIST_HEIGHT_BETWEEN_ROWS, METRIC_ROW_Y_HEIGHT)
    m.create_measurement("date", 1, 0, DATE_ROW_Y_POS, 0, [])
    m.create_measurement("weight", 2, 0, 0, 0, [60, 99])
    m.create_measurement("body_mass_index", 2, 1, 0, 0, [20, 35])
    m.create_measurement("body_fat_percentage", 2, 0, 0, 1, [0, 100])
    m.create_measurement("water_percentage", 2, 1, 0, 1, [0, 100])
    m.create_measurement("muscle_mass_percentage", 2, 0, 0, 2, [0, 100])
    m.create_measurement("bone_mass_percentage", 2, 1, 0, 2, [0, 100])
    m.create_measurement("basal_metabolic_rate", 2, 0, 0, 3, [1400, 2000])
    m.create_measurement("visceral_fat", 2, 1, 0, 3, [1, 20])
    m.create_measurement("lean_body_mass", 2, 0, 0, 4, [50, 80])
    m.create_measurement("body_fat_mass", 2, 1, 0, 4, [10, 30])
    m.create_measurement("bone_mass", 2, 0, 0, 5, [2, 4])
    m.create_measurement("muscle_mass", 2, 1, 0, 5, [40, 80])
    m.create_measurement("body_age", 2, 0, 0, 6, [20, 35])
    m.create_measurement("protein_percentage", 2, 1, 0, 6, [10, 50])
    # crop loop
    for key, value in m.measurement_list.items():
        crop_img = image[value.y:value.y+value.h, value.x:value.x+value.w]
        # feed images to tesseract
        if key != "date":
            proc_image = process_metric_unit_for_tesseract(crop_img)
            text = pytesseract.image_to_string(
                proc_image, config="-c tessedit_char_whitelist=.0123456789 --psm 6").rstrip()
            text = process_metric_output(text, value.reasonable_range, logger)

        else:
            proc_image = process_date_for_tesseract(crop_img)
            text = pytesseract.image_to_string(proc_image).rstrip()
            text = process_date_output(text, logger)

        value.value = text
    return m


def main():
    logger = logging.getLogger('magic-extractor')
    logger.setLevel(logging.DEBUG)
    fh = logging.FileHandler('logs/magic-extractor.log')
    fh.setLevel(logging.DEBUG)
    logger.addHandler(fh)
    image = cv2.imread(sys.argv[1])
    h, w, c = image.shape
    if w != SCREENSHOT_WIDTH or h != SCREENSHOT_HEIGHT:
        logger.error("error: screenshot not the correct resolution")
        exit()

    measurements = process_image(image, logger)
    logger.info(measurements)
    print(measurements.toJSON())


main()
