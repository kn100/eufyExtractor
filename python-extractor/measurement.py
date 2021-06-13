class Measurement:
    def __init__(self, x, y, w, h, reasonable_range):
        self.x = x
        self.y = y
        self.w = w
        self.h = h
        self.val = "unset"
        self.reasonable_range = reasonable_range

    def set_value(self, val):
        self.val = val
