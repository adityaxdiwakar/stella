from bs4 import BeautifulSoup
from dotenv import load_dotenv
import datetime
import requests
import os

load_dotenv() #grab env variables from config

prefix = os.getenv("BOT_PREFIX")
is_dev = bool(int(os.getenv("IS_DEV")))

from commands import ping
from commands import charts
from commands import futcharts
from commands import fxcharts
from commands import refs
from commands import ng_all
from commands import ng_rep
from commands import earnings
from commands import eightball
from commands import evalmod
from commands import custom_futures

import copy
import discord
import time
import asyncio
import ast

def exec_then_eval(code):
    block = ast.parse(code, mode='exec')

    # assumes last node is an expression
    last = ast.Expression(block.body.pop().value)

    _globals, _locals = {}, {}
    exec(compile(block, '<string>', mode='exec'), _globals, _locals)
    return eval(compile(last, '<string>', mode='eval'), _globals, _locals)

async def status():
    counter = 0
    links = ["https://www.investing.com/indices/us-spx-500-futures", "https://www.investing.com/indices/us-spx-500-futures", "https://www.investing.com/indices/us-spx-500-futures", "https://www.investing.com/indices/us-spx-500-futures"]
    tickers = ["ES", "ES", "ES", "ES"]
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


# {prefix: component}
module_links = {
    "ping": ping.main,
    "cf": custom_futures.main,
    "c": charts.main,
    "mc": charts.multi,
    "f": futcharts.main,
    "x": fxcharts.main,
    "ngall": ng_all.main,
    "rc": ng_rep.custom,
    "anom": ng_rep.all_anom,
    "r": ng_rep.main,
    "earnings": earnings.company,
    "addtag": refs.add_ref,
    "8ball": eightball.main,
    "showtags": refs.show_tags,
    "tag": refs.use_tag,
    "eval": evalmod.main
}

class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)
        channel = ctx.get_channel(636986005773352980)
        dev_msg = "I am currently running in the **production** environment."
        if is_dev:
            dev_msg = "I am currently running in a **canary development** environment."
        await channel.send(f"Stella has been rebooted. The current time is {datetime.datetime.now().strftime('%H:%M:%S on %b %d')}. {dev_msg}")
        ctx.loop.create_task(status())

    async def on_message(self, message):
        message.content.split(" ")[0] = message.content.split(" ")[0].lower()
        if message.content.startswith(f"{prefix}"):
            if message.content[len(prefix)] == " ":
                char_array = list(message.content)
                del char_array[len(prefix)]
                message.content = ""
                for char in char_array:
                    message.content += char

        for mod in module_links:
            if message.content.startswith(prefix + mod):
                await module_links[mod](message, canary=is_dev)
                break

        

ctx = Stella()
ctx.run(os.getenv("BOT_TOKEN"))
