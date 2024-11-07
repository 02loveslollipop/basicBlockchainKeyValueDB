import requests

class NoPayloadException(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)
        
class HostNotAvailableException(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)
        
class ServerCoulNotDecodeJSON(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)
        
class NoConsensusException(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)
    
class InternalServerError(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)

class blockchainDB:
    def __init__(self, host: str, port: int):
        self.host = host
        self.port = port
        
    def add(self, key: str = None, value: str = None, data: dict = None) -> None:
        if data is None: # If no data (dictionary) is provided
            if key is None or value is None: # Check if key and value are provided
                raise NoPayloadException('Please provide a key and a value or a valid dictionary') # If either key or value nor dictionary is provided, raise an exception
            payload = { # If key and value are provided, create a dictionary from them
                "key": key,
                "value": value
            }
        elif data.keys() == 0: # If an empty dictionary is provided
            raise NoPayloadException('Please provide a key and a value or a valid dictionary')
        elif data.keys() > 1: # If more than one key-value pair is provided
            raise NoPayloadException('Please provide only one key-value pair')
        else:
            payload = data # If a valid dictionary is provided, use it as the payload
            
        try:
            response = requests.post(f'http://{self.host}:{self.port}/append', json=payload) # Send a POST request to the server
        except requests.exceptions.ConnectionError as e:
            raise HostNotAvailableException('Could not connect to the server')
        if response.status_code == 201:
            return
        if response.status_code == 400:
            raise ServerCoulNotDecodeJSON('Invalid payload, please provide a valid JSON string')
        if response.status_code == 401:
            raise NoConsensusException('No consensus, please try again later')
        if response.status_code == 500:
            raise InternalServerError(f'Internal server error: {response.text}')
    
    def fetch(self) -> dict:
        try:
            response = requests.get(f'http://{self.host}:{self.port}/chain') # Send a GET request to the server
        except requests.exceptions.ConnectionError:
            raise HostNotAvailableException('Could not connect to the server')
        if response.status_code == 200:
            return response.json()
        if response.status_code == 500:
            raise InternalServerError(f'Internal server error: {response.text}')