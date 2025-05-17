import requests
import json

def main():
    url = "http://localhost:8081/semantic-search"
    
    print("Semantic Search Client")
    print("Type 'exit' to quit.\n")
    
    while True:
        user_input = input("Enter your search query: ")
        if user_input.lower() == "exit":
            print("Exiting...")
            break

        payload = {
            "search": user_input,
            "max_n": 5
        }
        
        try:
            response = requests.post(url, json=payload)
            if response.status_code == 200:
                # print("Response:", json.dumps(response.json(), indent=2))
                print(response.json()["summary"])
            else:
                print(f"Error: {response.status_code} - {response.text}")
        except requests.exceptions.RequestException as e:
            print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()