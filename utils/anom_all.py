import plotly.graph_objects as go
from plotly.subplots import make_subplots

import requests
import os
import time

def main(report_type):
    try:
        fig = go.Figure()
        r_data = requests.get(f"https://api.adi.wtf/weather/reports/{report_type}").json()
        climo_data = requests.get("https://api.adi.wtf/weather/historical/climo").json()["hdd"]
        obs_data = requests.get("https://api.adi.wtf/weather/historical/observations").json()["hdd"]

        key = list(r_data.keys())[-1]
        hdd_vals = r_data[key]["hdd"]
        hdd_diffs = {x:float(hdd_vals[x]) - float(climo_data[x]) for x in list(hdd_vals.keys()) if hdd_vals[x] != None and x in climo_data}
        fig.add_trace(go.Scatter(x=list(hdd_diffs.keys()), y=list(hdd_diffs.values()), name=key, line_shape='spline'))

        obs_data = {x:obs_data[x] for x in sorted(obs_data.keys())}
        obs_diffs = {x:float(obs_data[x]) - float(climo_data[x]) for x in list(obs_data.keys()) if x in climo_data and obs_data[x] != None}
        fig.add_trace(go.Scatter(x=list(obs_diffs.keys()), y=list(obs_diffs.values()), name="Observations", line_shape='spline'))

        fig.update_layout(shapes=[go.layout.Shape(type='line', xref='paper', x0=0, x1=1, y0=0, y1=0, line=dict(color='white', dash='dash'))])
        fig.update_layout(template="plotly_dark", title=f"{report_type.upper()} Report Chart", width=1440)

        rn = time.time()
        fig.write_image(f"{os.getenv('SYS_FILE_LOCATION')}/{rn}.jpeg")

        return f"{os.getenv('FILE_LOCATION')}/{rn}.jpeg"
    except Exception as e:
        return f"```\n{e}\n```"
