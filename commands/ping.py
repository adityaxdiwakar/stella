from datetime import datetime
import time

async def main(message):
    msg = await message.channel.send(f":ping_pong: Calculating ping...")
    roundtrip = round((msg.created_at.timestamp() - message.created_at.timestamp()) * 1000)
    await msg.edit(content=f":ping_pong: WS roundtrip complete in {roundtrip}ms!")