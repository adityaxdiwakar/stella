from bs4 import BeautifulSoup
from dotenv import load_dotenv
import requests
import os

load_dotenv() #grab env variables from config

prefix = os.getenv("BOT_PREFIX")
is_dev = bool(int(os.getenv("IS_DEV")))

from commands import ping
from commands import charts
from commands import futcharts
from commands import fxcharts

import copy
import discord
	import time
import asyncio

async def status():
    counter = 0
    links = ["https://www.investing.com/commodities/gold", "https://www.investing.com/indices/us-spx-500-futures", "https://www.investing.com/indices/us-spx-500-futures", "https://www.investing.com/indices/us-spx-500-futures"]
    tickers = ["GC", "ES", "ES", "ES"]
    while True:
        r = requests.get(links[counter % len(links)],  headers={'User-Agent': 'Mozilla/5.0'})
        soup = BeautifulSoup(r.content, 'html.parser')
        last_price_obj = soup.find(id="last_last")
        prices = [str(div.string).strip() for div in last_price_obj.parent]
        prices = [div for div in prices if div != ""]
        last_price = prices[0]
        last_price = float(prices[0].replace(",", ""))
        price_change = prices[1]
        price_percent = prices[2][1:]
        ticker = tickers[counter % len(links)]
        activity = discord.Activity(name=f"{ticker}: {last_price} ({price_change} {price_percent})", type=0)
        await ctx.change_presence(activity=activity)
        counter += 1
        await asyncio.sleep(15)

class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)
        ctx.loop.create_task(status())

    async def on_message(self, message):
        message.content = message.content.lower()
        if message.content.startswith(f"{prefix}"):
            if message.content[len(prefix)] == " ":
                char_array = list(message.content)
                del char_array[len(prefix)]
                message.content = ""
                for char in char_array:
                    message.content += char

            if message.content.startswith(f"{prefix}ping"):
                await ping.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}c"):
                await charts.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}f"):
                await futcharts.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}x"):
                await fxcharts.main(message, canary=is_dev)

ctx = Stella()
ctx.run(os.getenv("BOT_TOKEN"))
