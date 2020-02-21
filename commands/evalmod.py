async def main(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    if message.author.id == 192696739981950976 or message.author.id == 513549665019363329:
        command = " ".join(message.content.split(" ")[1:])
        print(command)
        try:
            resp = eval(command)
        except Exception as e:
            resp = e

        msg = f"```{resp}```"
        await message.channel.send(prefix + msg)
    else:
        await message.channel.send(prefix + ":lock: You do not have permission to evaluate in runtime!")
