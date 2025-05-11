import requests

# List of URLs to post
urls = [
    "https://www.bbcgoodfood.com/recipes/one-pot-paneer-curry-pie",
    "https://www.bbcgoodfood.com/recipes/salt-and-pepper-chips",
    "https://www.bbcgoodfood.com/recipes/birria-tacos",
    "https://www.bbcgoodfood.com/recipes/microwave-cheats-paella",
    "https://www.bbcgoodfood.com/recipes/paneer-korma",
    "https://www.bbcgoodfood.com/recipes/easy-butter-chicken",
    "https://www.bbcgoodfood.com/recipes/salt-pepper-tofu",
    "https://www.bbcgoodfood.com/recipes/chorizo-pea-risotto",
    "https://www.bbcgoodfood.com/recipes/gnocchi-fish-pie",
    # Add more URLs as needed
]

# Endpoint to send POST requests
endpoint = "http://localhost:8081/import-url"

# Loop through each URL and send a POST request
for url in urls:
    payload = {"url": url}
    try:
        response = requests.post(endpoint, json=payload)
        if response.status_code == 200:
            print(f"Successfully posted URL: {url}")
        else:
            print(f"Failed to post URL: {url}. Status Code: {response.status_code}, Response: {response.text}")
    except requests.RequestException as e:
        print(f"An error occurred while posting URL: {url}. Error: {e}")