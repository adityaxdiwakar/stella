reaction_config = {
    632602390985703424: { # verification channel 
        "💰": 632600011280089108 # moneybag and relevant role
    },
    691130032051060766: {
        "🕐": 691116986272579614,
        "🛶": 691117034699751464,
        "📆": 691117056313262111,
        "🖱️": 691117080472453161,
        "🌉": 691117127389675561,
        "💵": 691117221556125789,
        "📈": 691117223179190333
    }
}

async def handler(ctx, payload, state):
    m_id = payload.message_id
    if m_id not in reaction_config:
        return

    e_name = payload.emoji.name
    if e_name not in reaction_config[m_id]:
        return

    g_id = payload.guild_id
    u_id = payload.user_id
    r_id = reaction_config[m_id][e_name]
    guild = ctx.get_guild(g_id)
    role = guild.get_role(r_id)
    member = guild.get_member(u_id)

    if state == "add":
        await member.add_roles(role, reason="Reaction request fulfillment.")
    else:
        await member.remove_roles(role, reason="Reaction request fulfillment.")
