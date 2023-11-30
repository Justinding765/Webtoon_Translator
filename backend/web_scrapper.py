import requests
from bs4 import BeautifulSoup

def scrape_to_file(url, output_file):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raises an HTTPError if the HTTP request returned an unsuccessful status code

        # Parse the HTML content
        soup = BeautifulSoup(response.text, 'html.parser')

        # Find the specific div by its ID
        def match_id(tag):
            return tag.name == 'div' and tag.get('id') in ['comic_view_area', 'readerarea']
        view_area_div = soup.find(match_id)

        # Write the specific div's content to a file
        if view_area_div:
            with open(output_file, 'w', encoding='utf-8') as file:
                file.write(str(view_area_div))
                print(f"Content of div 'comic_view_area' written to {output_file}")
        else:
            print("Div with id 'comic_view_area' not found.")

    except requests.RequestException as e:
        print(f"Error during requests to {url}: {str(e)}")

# Example usage
scrape_to_file("https://cosmic-scans.com/eleceed-raws-chapter-273/", "output.html")
