import plotly.graph_objects as go
from plotly.subplots import make_subplots

import requests
import os
import time

def main(feed):
    try:
        fig = go.Figure()

        r_data = {}
        for rt in feed:
            r_data.update({
                rt.split("@")[0]: requests.get(f"https://api.adi.wtf/weather/reports/{rt.split('@')[0]}").json()
            })

        for rt in feed:
            report = rt.split("@")[0]
            date = rt.split("@")[1].replace('#', ' ')

            vals = r_data[report][date]["hdd"]
            fig.add_trace(go.Scatter(x=list(vals.keys()), y=list(vals.values()), name=date, line_shape='spline'))

        #fig.update_layout(shapes=[go.layout.Shape(type='line', xref='paper', x0=0, x1=1, y0=0, y1=0, line=dict(color='white', dash='dash'))])

        fig.update_layout(template="plotly_dark", title=f"Custom Report Chart", width=1440)

        rn = time.time()
        fig.write_image(f"{os.getenv('SYS_FILE_LOCATION')}/{rn}.jpeg")

        return f"{os.getenv('FILE_LOCATION')}/{rn}.jpeg"
    except Exception as e:
        return f"```{e}```"
