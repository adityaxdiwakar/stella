import plotly.graph_objects as go

import requests
import time
import json

async def main(message, canary=False):
    try:
        premsg = ["**[Canary]** ", " "][not canary]

        msg = await message.channel.send(premsg + "Okay, standby while the chart is generated.")
        
        r = requests.get("http://vixcentral.com/ajax_update", headers={
            "X-Requested-With": "XMLHttpRequest",
        })

        if r.status_code != 200:
            raise Exception("Error, handle in exception handler.")

        r_json = r.json()

        contract_months = r_json[0]
        futures_prices = r_json[2]
        index_prices = r_json[8]
        index_text = [None for x in range(len(index_prices))]
        index_text[-1] = index_prices[0]

        fig = go.Figure()

        fig.add_trace(go.Scatter(x=contract_months, 
                                y=futures_prices, 
                                name="Futures",
                                line_shape='spline', 
                                mode='lines+markers+text',
                                text=[str(x) for x in futures_prices],
                                textposition="top center"))
        
        fig.add_trace(go.Scatter(x=contract_months, 
                                y=index_prices, 
                                name="VIX Index", 
                                line=dict(
                                    color="green",
                                    dash='dash'
                                ),
                                line_shape='spline', 
                                mode='lines+markers+text',
                                text=index_text,
                                textposition="bottom center"))

        fig.update_traces(textfont_size=18)
        fig.update_layout(template="plotly_dark", title=f"VIX Calendar Curve for {r_json[1][0].split(' ')[0]}", width=1440)

        rn = time.time()
        fig.write_image(f"bin/{rn}.jpeg")
        link = f"https://img.adi.wtf/ca/{rn}.jpeg"
        # fig.write_image(f"/var/www/html/u/ca/{rn}.jpeg")

        await msg.edit(content=premsg + f"Here you are: {link}")
        
    except Exception as e:
        print(e)
        await msg.edit(content=premsg + "Unfortunately, an error has occured. Developers have been notified of the issues.")