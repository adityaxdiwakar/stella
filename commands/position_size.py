async def calculator(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    msg_split = message.content.split(" ")
    
    if len(msg_split) < 2:
        await message.channel.send(premsg + "Sorry, please enter what the value of VIX is!")
        return

    try:
        vix_value = float(msg_split[1])
    except:
        await message.channel.send(premsg + "Please enter a number for VIX, I cannot interpret your message")
        return

    t_1 = -0.000007326 * (vix_value**2) 
    t_2 = 0.002381 * (vix_value)
    t_3 = t_1 - t_2
    t_4 = t_3 + 0.4374
    t_4 = round(t_4, 3)

    await message.channel.send(premsg + f"Alright, the position proportion for today given VIX is {t_4}")