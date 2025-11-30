#!/bin/bash
set -e

INSTALL_PATH=/opt/snitchdns

cd $INSTALL_PATH

# Create test admin user if not exists
echo "Creating test user..."
./venv.sh flask users add \
    --username testadmin \
    --password '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5LL2l0Q7o0xei' \
    --full_name "Test Administrator" \
    --email "test@example.com" \
    --active yes \
    --admin yes \
    --auth local \
    2>/dev/null || echo "User already exists or creation failed - continuing..."

# Get or create API key for test user
echo "Setting up API key..."
USER_ID=$(./venv.sh flask users list | grep testadmin | awk '{print $1}' | head -1)

if [ ! -z "$USER_ID" ]; then
    # Create API key via Python script
    ./venv.sh python3 << 'PYTHON_SCRIPT'
from app import create_app
from app.lib.base.provider import Provider

app = create_app()
with app.app_context():
    provider = Provider()
    api = provider.api()
    users = provider.users()

    user = users.find_user_login('testadmin')
    if user:
        # Check if API key exists
        keys = api.all(user_id=user.id)
        if len(keys) == 0:
            # Create new API key
            apikey = api.add(user.id, 'testkey')
            print(f"API_KEY={apikey.apikey}")
        else:
            print(f"API_KEY={keys[0].apikey}")
PYTHON_SCRIPT
fi > /tmp/apikey.txt

# Export API key for other processes if needed
if [ -f /tmp/apikey.txt ]; then
    export $(cat /tmp/apikey.txt | xargs)
    echo "API Key configured: ${API_KEY:0:10}..."
fi

# Start services
echo "Starting cron..."
service cron start

echo "Starting SnitchDNS daemon..."
./venv.sh flask snitch_start

echo "Starting Flask web application on port 80..."
exec ./venv.sh flask run --host 0.0.0.0 --port 80
