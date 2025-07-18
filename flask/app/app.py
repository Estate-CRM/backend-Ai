from flask import Flask
from app.routes import setup_routes  # Import the function from routes.py

app = Flask(__name__)
setup_routes(app)  # Register your routes here
