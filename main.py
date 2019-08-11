from dotenv import load_dotenv
import os

load_dotenv() #grab env variables from config

from commands import ping
from commands import charts

import copy
import discord
import time

class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)

    async def on_message(self, message):
        if message.content.startswith("?"):
                if message.content.startswith("?ping"):
                    await ping.main(message)

                elif message.content.startswith("?c"):
                    await charts.main(message)
                    
        
ctx = Stella()
ctx.run(os.getenv("BOT_TOKEN"))