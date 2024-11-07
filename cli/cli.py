import requests
import click
import json

def parse_cli_payload(ctx, param, value):
    if value is None:
        return None
    try:
        return json.loads("{" + value + "}")
    except json.JSONDecodeError:
        raise click.BadParameter('Please provide a valid JSON string. Example: --data "key":"value"')

@click.group()
def cli():
    """A simple CLI for interacting with the blockchain"""
    pass

@cli.command()
@click.argument('data', callback=parse_cli_payload)
@click.option('--url', help='Server URL', default='localhost')
@click.option('--port', help='Server port', default='8080')
def add(data: dict, url, port):
    """Add a key value pair to the blockchain"""
    try:
        key = list(data.keys())[0]
        value = list(data.values())[0]
        payload = {
            "key": key,
            "value": value
        }
        click.echo(payload)
        response = requests.post(f'http://{url}:{port}/append', json=payload)
    except requests.exceptions.ConnectionError:
        click.echo('Could not connect to the server')
        return
    if response.status_code == 201:
        click.echo('Data added successfully')
        return
    if response.status_code == 400:
        click.echo('Invalid payload, please provide a valid JSON string')
        return
    if response.status_code == 401:
        click.echo('No consensus, please try again later')
        return
    if response.status_code == 500:
        click.echo(f'Internal server error: {response.text}')
        return
    
@cli.command()
@click.option('--url', help='Server URL', default='localhost')
@click.option('--port', help='Server port', default='8080')
def rawfetch(url, port):
    """Fetch the blockchain and return it as is"""
    try:
        response = requests.get(f'http://{url}:{port}/chain')
    except requests.exceptions.ConnectionError:
        click.echo('Could not connect to the server')
        return
    if response.status_code == 200:
        click.echo(response.json())
        return
    if response.status_code == 500:
        click.echo(f'Internal server error: {response.text}')
        return

@cli.command()
@click.option('--url', help='Server URL', default='localhost')
@click.option('--port', help='Server port', default='8080')
def fetch(url, port):
    """Fetch the blockchain and return it as a list of dictionaries, excluding the Genesis block"""
    try:
        response = requests.get(f'http://{url}:{port}/chain')
    except requests.exceptions.ConnectionError:
        click.echo('Could not connect to the server')
        return
    if response.status_code == 500:
        click.echo(f'Internal server error: {response.text}')
        return
    if response.status_code != 200:
        click.echo('An error occurred')
        return
    
    chain = response.json()
    result = []
    for block in chain:
        data = block['data']
        if data['key'] == 'Genesis' and data['value'] == 'Genesis block':
            continue
        else:
            result.append({data['key']: data['value']})
    click.echo(result)

if __name__ == '__main__':
    cli()