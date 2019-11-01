
import plotly.graph_objects as go
from plotly.subplots import make_subplots

import requests
import os
import time

def main(report_type):
    try:
        fig = go.Figure()
        r_data = requests.get(f"https://api.adi.wtf/weather/reports/{report_type}").json()
        climo_data = requests.get(f"https://api.adi.wtf/weather/historical/climo").json()["hdd"]

        inds = [-2, -1]
        if report_type.startswith("g"):
            inds = [-3, -1]

        sums = []
        climo_sum = 0
        for ind in inds:
            key = list(r_data.keys())[ind]
            hdd_vals = r_data[key]["hdd"]
            hdd_diffs = {x:float(hdd_vals[x]) - float(climo_data[x]) for x in list(hdd_vals.keys()) if hdd_vals[x] != None}
            fig.add_trace(go.Scatter(x=list(hdd_diffs.keys()), y=list(hdd_diffs.values()), name=key, line_shape='spline'))
            sums.append(sum([float(hdd_vals[x]) for x in list(hdd_vals.keys()) if hdd_vals[x] != None]))
            if ind == -1:
                climo_sum = sum([float(climo_data[x]) for x in list(hdd_vals.keys())])

        sum_diff = round(sums[1] - sums[0], 2)
        climo_diff = round(sums[1] - climo_sum, 2)

        fig.update_layout(shapes=[go.layout.Shape(type='line', xref='paper', x0=0, x1=1, y0=0, y1=0, line=dict(color='white', dash='dash'))])

        fig.update_layout(template="plotly_dark", title=f"{report_type.upper()} Report Chart", width=1440)

        rn = time.time()
        fig.write_image(f"{os.getenv('SYS_FILE_LOCATION')}/{rn}.jpeg")

        return [f"{os.getenv('FILE_LOCATION')}/{rn}.jpeg", sum_diff, climo_diff]
    except Exception as e:
        return [f"```\n{e}\n```", None, None]
