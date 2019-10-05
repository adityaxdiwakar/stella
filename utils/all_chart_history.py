import plotly.graph_objects as go
from plotly.subplots import make_subplots

import requests
import os
import time

def main():
    try:
        positions = [(1, 1), (1, 2), (2, 1), (2, 2)]
        reports = ["gfs", "gfs-ensemble", "ecmwf-ensemble", "ecmwf"]
        fig = make_subplots(rows=2, cols=2, subplot_titles=reports, shared_yaxes=True, vertical_spacing=0.02)
        for y in range(len(reports)):
            r = requests.get(f"https://api.adi.wtf/weather/reports/{reports[y]}").json()
            for x in r:
                fig.add_trace(
                    go.Scatter(x=list(r[x]["hdd"].keys()), y=list(r[x]["hdd"].values()), name=x, line_shape='spline'),
                    row= positions[y][0], col = positions[y][1]
                )    

        fig.update_layout(template="plotly_dark", title="All Report Charts", width=3840, height=2160)

        rn = time.time()
        fig.write_image(f"{os.getenv('SYS_FILE_LOCATION')}/{rn}.jpeg")

        return f"{os.getenv('FILE_LOCATION')}/{rn}.jpeg"
    except Exception as e:
        return f"```\n{e}\n```"