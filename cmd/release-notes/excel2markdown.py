import pandas as pd
import argparse
import re

parser = argparse.ArgumentParser(description='Convert an Excel file to Markdown')
parser.add_argument('input_file', type=str, help='path to input Excel file')
parser.add_argument('output_file', type=str, help='path to output Markdown file')
args = parser.parse_args()

df = pd.read_excel(args.input_file)

result = {}

def get_release_note_and_link(note):
    # Example:
    #   Add the Observability Center. ([kse-console#195](https://github.com/kubesphere/kse-console/pull/195), [@harrisonliu5](https://github.com/harrisonliu5))
    pr_link = note["pr_url"]
    match = re.match(r"https://github.com/kubesphere/(.*)/pull/(\d+)", pr_link)
    pr_link_text = "%s#%s" % (match.group(1), match.group(2))
    release_note = note["release_note"]
    author = note["author"]
    author_link = "https://github.com/%s" % author
    return "%s ([%s](%s), [@%s](%s))" % (release_note, pr_link_text, pr_link, author, author_link)

for _, row in df.iterrows():
    if isinstance(row['Release Note'], str):
        module = result.setdefault(row['Module'], {})
        type = module.setdefault(row['Type'], [])
        if isinstance(row['pr_url'], str):
            type.append({
                "release_note": row['Release Note'],
                "pr_url": row['pr_url'],
                "author": row['author'],
            })

md = ""
for module, types in result.items():

    md += "## %s\n\n" % module
    for type, notes in types.items():
        md += "### %s\n\n" % type
        for note in notes:
            md += "- %s\n" % get_release_note_and_link(note)
        md += '\n'
    md += '\n'


# 写入Markdown文件
with open(args.output_file, 'w') as f:
    f.write(md)

print('Conversion complete!')