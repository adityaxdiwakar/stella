from bs4 import BeautifulSoup
import requests
import discord

async def company(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    msg = await message.channel.send(premsg + "Loading...")

    if len(message.content.split(" ")) < 2:
        await msg.edit(content=premsg + "Sorry, please enter the ticker you would like to lookup earnings for!")
        return

    ticker = message.content.split(" ")[1]

    r = requests.get(f"https://earningswhispers.com/stocks/{ticker}")

    soup = BeautifulSoup(r.content, 'html.parser')
    whisper_box = soup.find(id="whisperbox")
    date_box = soup.find(id="datebox")
    try:
        company_name = soup.find(id="compname").text
    except AttributeError:
        await msg.edit(content=premsg + "Sorry, that company does not exist in the earnings database. Try again?")
        return

    print(company_name)

    date_box_children = [child.text for child in date_box.children if child.text != ""]
    print(date_box_children)

    whisper_box_children = [child.text for child in whisper_box.children if child.text != ""]
    print(whisper_box_children)

    embed = discord.Embed(title=ticker.upper(), description=f"{company_name} Earnings ({whisper_box_children[2]})")
    embed.add_field(name="Reporting Date", value=f"{date_box_children[1]} ({date_box_children[0]})", inline=True)
    embed.add_field(name="Reporting Time", value=date_box_children[2], inline=True)
    embed.add_field(name="EW EPS", value=date_box_children[1], inline=True)
    print(whisper_box_children[3].split(" "))
    embed.add_field(name="WS EPS", value=whisper_box_children[3].split(" ")[2], inline=True)
    embed.add_field(name="WS Revenue", value=" ".join(whisper_box_children[4].split(" ")[1:3]), inline=True)


    await msg.edit(content=premsg, embed=embed)
