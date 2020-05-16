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
from commands import position_size
from commands import tdcommands
from commands import vixcentral

from utils import reactions

import copy
import discord
import time
import asyncio
import ast

# {prefix: component}
module_links = {
    "deltag": refs.rm_tag,
    "ping": ping.main,
    "c": charts.main,
    "mc": charts.multi,
    "fun": tdcommands.fundamentals,
    "div": tdcommands.dividends,
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
    "eval": evalmod.main,
    "pos": position_size.calculator,
    "vixc": vixcentral.main
}

async def update_price():
    await ctx.wait_until_ready()
    ws_channel = ctx.get_channel(703080609358020608)
    em_channel = ctx.get_channel(709860290694742098)
    while True:
        try:
            r = requests.get("https://md.adi.wtf/recent/")
            price = r.json()["payload"]["trade"]["price"]
            settlement = r.json()["payload"]["session_prices"]["settlement"]
            n_percentage = round(100 * (price - settlement) / 2826.5, 2)
            percentage = str(n_percentage) + "%"
            if n_percentage > 0:
                percentage = "+" + percentage
            message = f"ES @ {price} ({percentage})"
            await ws_channel.edit(name=message)
            await em_channel.edit(name=message)
        except:
            pass
        await asyncio.sleep(12) # task runs every 60 seconds

class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)

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

    async def on_raw_reaction_add(self, payload):
        await reactions.handler(self, payload, "add")
    
    async def on_raw_reaction_remove(self, payload):
        await reactions.handler(self, payload, "remove")

        

ctx = Stella()
ctx.loop.create_task(update_price())
ctx.run(os.getenv("BOT_TOKEN"))
