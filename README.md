Stella | Discord Bot
----

Stella is a Discord Bot intended to provide access to financial data through open and closed source tools provided by third-parties and Aditya Diwakar as well. This bot is intended to be used in the ``Wall St. Community`` server but has been introduced in numerous servers since then. This bot is written in Golang and is currently on-going a full migration from the previous Python version.

## Installation
The following Shell script can be used to update Stella, assuming that Stella is being run using ``systemd``:
```sh
rm stella
curl -s https://api.github.com/repos/adityaxdiwakar/stella/releases/latest \
| grep "browser_download_url.*" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi -
chmod +x stella
systemctl restart stella-v2
```

## Support

Support for this bot is provided, please open a GitHub issue or contact me through Discord (``aditya#0001``) or email (``aditya@diwakar.io``).

## Contribution

We're currently closed to taking any requests, if you make a PR, it'll most likely be denied due to our policies.

