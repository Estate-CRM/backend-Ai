from flask import Flask
from .consumer import start_consumer_thread  # ✅ relative import

app = Flask(__name__)

# Optional: from .routes import setup_routes
# setup_routes(app)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
