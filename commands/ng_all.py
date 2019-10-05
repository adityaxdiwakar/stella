from utils import all_chart_history as ach
from datetime import datetime
import time

async def main(message, canary=False):
    premsgs = ["**[Canary]** ", " "]
    premsg = premsgs[not canary]
    msg = await message.channel.send(premsg + "Generating chart, stand by.")
    chart_link = ach.main()
    await msg.edit(content=premsg + f"Here you go: {chart_link}")