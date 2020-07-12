rm stella
curl -s https://api.github.com/repos/adityaxdiwakar/stella/releases/latest \
    | grep "browser_download_url.*" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qi -
chmod +x stella
systemctl restart stella-v2
