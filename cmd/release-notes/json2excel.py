import click
import pandas as pd
import json

@click.command()
@click.argument('json_file', type=click.File('r'))
@click.argument('output_file', default='output.xlsx', type=str)
def json2excel(json_file, output_file):
    """
    Converts a JSON file to an Excel file

    JSON_FILE: Path to the JSON file to be converted
    OUTPUT_FILE: Path to the output Excel file (default: output.xlsx)
    """
    # Load data from the JSON file
    data = json.load(json_file)

    # Create a DataFrame object from the data
    df = pd.DataFrame.from_dict(data, orient='index')

    # Output the DataFrame to an Excel file
    df.to_excel(output_file, index=False)

    click.echo("Excel file generated: " + output_file)

if __name__ == '__main__':
    json2excel()