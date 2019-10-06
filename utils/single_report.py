import plotly.graph_objects as go
from plotly.subplots import make_subplots

import requests
import os
import time

def main(report_type):
    try:
        fig = go.Figure()
        r_data = requests.get(f"https://api.adi.wtf/weather/reports/{report_type}").json()
        climo_data = requests.get(f"https://api.adi.wtf/weather/historical/climo").json()

        
        fig.add_trace(go.Scatter(x=list(r[x]["hdd"].keys()), y=list(r[x]["hdd"].values()), name=x, line_shape='spline'))    

        fig.update_layout(template="plotly_dark", title="All Report Charts", width=3840, height=2160)

        rn = time.time()
        fig.write_image(f"{os.getenv('SYS_FILE_LOCATION')}/{rn}.jpeg")

        return f"{os.getenv('FILE_LOCATION')}/{rn}.jpeg"
    except Exception as e:
        return f"```\n{e}\n```"