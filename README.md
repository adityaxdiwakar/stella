# Stella | Discord Bot
----

Stella is a Discord Bot intended for the Stock Market Live Chat, it uses a combination of charting and statistical sources in order to provide quick and referencable data to the members of the server. For all questions and inquiries, feel free to contact ``@aditya#1337`` on Discord or email me ``aditya@diwakar.io``.

## Documentation

### Ping command

Our ping command is quite simple, run it using ``?ping``. You will be returned a simple response with the roundtrip time for the WS/API timing. This is extraneous and used to measure the health of the Bot/Discord servers.

### Charting commands

The charting command is encompassed in ``?c``.

#### Level System

The bot is equipped with a simple "level" system for the charting command:

0. 1 Minute Intraday 
1. 3 Minute Intraday
2. 5 Minute Intraday
3. 15 Minute Intraday
4. 30 Minute Intraday
5. Daily
6. Weekly
7. Monthly

These numbers can be appended to the command above, for example: ``?c2`` will bring up a 5 minute intraday chart.

By default, ``?c`` brings up a daily chart.

#### Usage

Use the command as: ``?c<level> <ticker>``. If the level is blank, the default will be used. Refer to the levels from above.

## Support

Support for this bot is provided, please open a GitHub issue or contact me through Discord (``aditya#1337``) or email (``aditya@diwakar.io``). 

## Contribution

We're currently closed to taking any requests, if you make a PR, it'll most likely be denied due to our policies.