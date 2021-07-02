import measurement
import json


class Measurements:
    def __init__(self, screen_width, first_metric_row_x_pos, dist_between_rows, y_height):
        self.screen_width = screen_width
        self.first_metric_row_x_pos = first_metric_row_x_pos
        self.dist_between_rows = dist_between_rows
        self.y_height = y_height
        self.measurement_list = {}

    def create_measurement(self, name, total_cols, col_num, y_pos_override, row_num, reasonable_range):
        w = int(self.screen_width / total_cols)
        h = self.y_height
        x = int((self.screen_width / total_cols) * col_num)
        y = 0
        if y_pos_override != 0:
            y = y_pos_override
        else:
            y = self.first_metric_row_x_pos + \
                (row_num * self.dist_between_rows)
        m = measurement.Measurement(x, y, w, h, reasonable_range)
        self.measurement_list[name] = m

    def getMeasurement(self, name):
        return self.measurement_list[name]

    def toJSON(self):
        jsons = {}
        jsons['measurements'] = []
        for key, value in self.measurement_list.items():
            if key != 'date':
                measurement = {}
                measurement["type"] = key
                measurement["value"] = value.value
                jsons["measurements"].append(measurement)
            else:
                jsons["date"] = int(value.value)

        return json.dumps(jsons)
