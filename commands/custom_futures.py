import zmq

import requests
import discord
import json
import time
import os

from datetime import datetime

import plotly.graph_objects as go
from plotly.subplots import make_subplots
context = zmq.Context()

prefix = os.getenv("BOT_PREFIX")

async def main(d_message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    chart_type = int(d_message.content[len(prefix) + 2])
    if chart_type == " ":
        chart_type = 1
    chart_type -= 1

    if chart_type < 0 or chart_type > 7:
        await d_message.channel.send(premsg + "Sorry, please pick between 1 and 8 for the chart length.")
        return

    split_message = d_message.content.split(" ")[1:]
 
    fig = go.Figure()

    ticker = split_message[0].upper()
    freq = ["m5", "m30", "h1", "h2", "d1", "w1", "w1", "n1"][chart_type]
    pd = ["d1", "w1", "m1", "m2", "y1", "y5", "y10", "y30"][chart_type]
    tick_format = ["%l:%M%p", "%a %l:%M%p\n", "%m/%d %l%p", "%m/%d %l%p",
                    "%b %d", "%b %d %Y", "%b %d %Y", "%b %Y"][chart_type]

    #  Socket to talk to server
    print("Connecting to market data hub...")
    socket = context.socket(zmq.PAIR)
    socket.connect("tcp://vps.adi.wtf:5556")
    print("Connected!")

    socket.send_string(f"/{ticker} {freq} {pd}")

    d_msg = await d_message.channel.send(premsg + "Creating chart, stand by.")

    message = socket.recv_string()
    print(message)
    message = socket.recv_json()
    print(message)

    data = message["snapshot"][0]["content"][0]["3"]
    data = data[:-1]

    try:
        times = [datetime.fromtimestamp(x["0"]/1000) for x in data][:-1]
        times = [x.strftime(tick_format) for x in times]
        opens = [x["1"] for x in data][:-1]
        highs = [x["2"] for x in data][:-1]
        lows = [x["3"] for x in data][:-1]
        closes = [x["4"] for x in data][:-1]
        volume = [x["5"] for x in data][:-1]

        fig.add_trace(go.Candlestick(x=times, open=opens, high=highs, low=lows, close=closes))

        fig.update_layout(shapes=[go.layout.Shape(type='line', xref='paper', x0=0, x1=1, y0=0, y1=0, line=dict(color='white', dash='dash'))])
        fig.update_layout(template="plotly_dark", title=f"Custom Futures Chart", width=4320, height=1440)
        fig.update_layout(yaxis=dict(range=[min(lows) - 0.001 * min(lows), max(highs) + 0.001 * max(highs)]), xaxis={'type': 'category'})
        fig.update_layout(xaxis_rangeslider_visible=False)
        fig.update_layout(font=dict(size=32))
        fig.update_xaxes(nticks=12)
        fig.update_layout(margin=go.layout.Margin(
            l=100,
            r=100,
            b=150,
            t=150,
            pad=20
        ))


        rn = time.time()
        fig.write_image(f"bin/{rn}.png")
        with open(f"bin/{rn}.png", "rb") as f:
            file = discord.File(f)
    except:
        await d_msg.edit(content=premsg + "Sorry, an error occured! Maybe that future contract doesn't exist?")
        return
    await d_message.channel.send(file=file)
    await d_msg.delete()