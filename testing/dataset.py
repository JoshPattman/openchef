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

urls = [
    "https://www.bbcgoodfood.com/recipes/amatriciana-chicken-traybake",
    "https://www.bbcgoodfood.com/recipes/oven-baked-risotto",
    "https://www.bbcgoodfood.com/recipes/chicken-pasta-bake",
    "https://www.bbcgoodfood.com/recipes/stickiest-ever-bbq-ribs-chive-dip",
    "https://www.bbcgoodfood.com/recipes/baked-ratatouille-goats-cheese",
    "https://www.bbcgoodfood.com/recipes/slow-roasted-mutton-shoulder-with-garlic-bean-mash-gremolata",
    "https://www.bbcgoodfood.com/recipes/pea-tarragon-cream-roast-chicken",
    "https://www.bbcgoodfood.com/recipes/roasted-tomato-basil-parmesan-quiche",
    "https://www.bbcgoodfood.com/recipes/fish-chip-pie",
    "https://www.bbcgoodfood.com/recipes/meat-potato-pie",
    "https://www.bbcgoodfood.com/recipes/baked-cauliflower-pizzaiola",
    "https://www.bbcgoodfood.com/recipes/creamy-salmon-leek-potato-traybake",
    "https://www.bbcgoodfood.com/recipes/flavour-bomb-roast-turkey-gravy",
    "https://www.bbcgoodfood.com/recipes/sausages-with-lemon-rosemary-roast-potatoes",
    "https://www.bbcgoodfood.com/recipes/stuffed-pasta-bake-bolognese",
    "https://www.bbcgoodfood.com/recipes/coronation-quiche",
    "https://www.bbcgoodfood.com/recipes/whole-baked-cauliflower-cheese",
    "https://www.bbcgoodfood.com/recipes/classic-lasagne-0",
    "https://www.bbcgoodfood.com/recipes/spiced-salmon-traybaked-sag-aloo",
    "https://www.bbcgoodfood.com/recipes/garlic-parmesan-breaded-chicken-with-quick-giardiniera",
    "https://www.bbcgoodfood.com/recipes/family-meals-easy-fish-pie-recipe",
    "https://www.bbcgoodfood.com/recipes/corned-beef-pie",
    "https://www.bbcgoodfood.com/recipes/pesto-sausage-traybake",
    "https://www.bbcgoodfood.com/recipes/chicken-taco-salad",
    "https://www.bbcgoodfood.com/recipes/orzo-tomato-soup",
    "https://www.bbcgoodfood.com/recipes/super-salad-wraps",
    "https://www.bbcgoodfood.com/recipes/buddha-bowl-salad",
    "https://www.bbcgoodfood.com/recipes/chicken-wraps",
    "https://www.bbcgoodfood.com/recipes/deli-style-stuffed-falafel-wrap",
    "https://www.bbcgoodfood.com/recipes/ham-mushroom-spinach-frittata",
    "https://www.bbcgoodfood.com/recipes/storecupboard-pasta-salad",
    "https://www.bbcgoodfood.com/recipes/caesar-pitta",
    "https://www.bbcgoodfood.com/recipes/curried-turkey-lettuce-wraps",
    "https://www.bbcgoodfood.com/recipes/tomato-black-bean-taco-salad",
    "https://www.bbcgoodfood.com/recipes/smoky-chickpea-salad",
    "https://www.bbcgoodfood.com/recipes/broad-bean-pea-ricotta-frittata",
    "https://www.bbcgoodfood.com/recipes/ponzu-tofu-poke-bowl",
    "https://www.bbcgoodfood.com/recipes/harissa-broccoli-flatbreads",
    "https://www.bbcgoodfood.com/recipes/chicken-tzatziki-wraps",
    "https://www.bbcgoodfood.com/recipes/monte-cristo-sandwich",
    "https://www.bbcgoodfood.com/recipes/air-fryer-ham-cheese-egg-bagel",
    "https://www.bbcgoodfood.com/recipes/spicy-bean-avocado-quesadillas",
    "https://www.bbcgoodfood.com/recipes/tuna-salad-sandwich",
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