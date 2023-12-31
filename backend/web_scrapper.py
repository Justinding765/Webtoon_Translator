import requests
from bs4 import BeautifulSoup
import sys



def match_id(tag):
            return tag.name == 'div' and tag.get('id') in ['comic_view_area', 'readerarea']
#Incase the specific id's aren't found
def general_scraping(soup):
    # Find all image tags
    images = soup.find_all('img')
    
    # Extract the 'src' attribute from each image
    image_urls = [str(img) for img in images]
    
    # Return the list of image URLs
    return '\n'.join(image_urls)

def scrape_page(url, output_file):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raises an HTTPError for unsuccessful status codes
        soup = BeautifulSoup(response.text, 'html.parser')
        view_area_div = soup.find(match_id)
        if view_area_div:
            with open(output_file, 'w', encoding='utf-8') as file:
                file.write(str(view_area_div))
                print(f"Content of div 'comic_view_area' written to {output_file}")
        else:
            # Fallback to general scraping (get all images)
            general_content = general_scraping(soup)
            with open(output_file, 'w', encoding='utf-8') as file:
                file.write(general_content)
                print(f"Image URLs written to {output_file}")

    except requests.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")
    except Exception as err:
        print(f"An error occurred: {err}")

if __name__ == '__main__':
    url = sys.argv[1]
    output_file = sys.argv[2]
    scrape_page(url, output_file)