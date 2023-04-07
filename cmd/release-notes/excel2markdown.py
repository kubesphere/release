import pandas as pd
import argparse

parser = argparse.ArgumentParser(description='Convert an Excel file to Markdown')
parser.add_argument('input_file', type=str, help='path to input Excel file')
parser.add_argument('output_file', type=str, help='path to output Markdown file')
args = parser.parse_args()

df = pd.read_excel(args.input_file)

result = {}

for _, row in df.iterrows():
    if isinstance(row['Release Note'], str):
        module = result.setdefault(row['Module'], {})
        type = module.setdefault(row['Type'], [])
        type.append(row['Release Note'])

md = ""
for module, types in result.items():

    md += "## %s \n\n" % module
    for type, notes in types.items():
        md += "### %s \n\n" % type
        for note in notes:
            md += "- %s \n" % note
        md += '\n'
    md += '\n'


# 写入Markdown文件
with open(args.output_file, 'w') as f:
    f.write(md)

print('Conversion complete!')