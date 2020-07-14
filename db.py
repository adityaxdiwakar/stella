import psycopg2

conn = psycopg2.connect("dbname=stella user=postgres password=password host=127.0.0.1 port=5432")
cur = conn.cursor()

import json
with open("bin/tags.json", "r") as f:
    data = json.load(f)

for tag in data:
    tag_name = tag
    tag_content = data[tag]
    sql = "INSERT INTO tags (id, content) VALUES (%s, %s)"
    cur.execute(sql, (tag_name, tag_content))
    conn.commit()

conn.close()
