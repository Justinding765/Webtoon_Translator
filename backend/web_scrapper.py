import requests
from bs4 import BeautifulSoup
import sys
import os
from google.cloud import storage



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

def scrape_page(url, SessionID):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raises an HTTPError for unsuccessful status codes
        soup = BeautifulSoup(response.text, 'html.parser')
        view_area_div = soup.find(match_id)

        script_dir = os.path.dirname(os.path.abspath(__file__))
        output_folder = os.path.join(script_dir, 'output')
        html_file = os.path.join(output_folder, f'{SessionID}.html')

        if view_area_div:
            with open(html_file, 'w', encoding='utf-8') as file:
                file.write(str(view_area_div))
                print(f"{html_file}")
        else:
            # Fallback to general scraping (get all images)
            general_content = general_scraping(soup)
            with open(html_file, 'w', encoding='utf-8') as file:
                file.write(general_content)
                print(f"Image URLs written to {html_file}")

    except requests.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")
    except Exception as err:
        print(f"An error occurred: {err}")
def create_bucket(sessionID):
    os.environ['GOOGLE_APPLICATION_CREDENTIALS'] = './config/credentials.json'

    # Create a Cloud Storage client.
    client = storage.Client()
     # Get the bucket with the name sessionID.
    bucket = client.bucket(sessionID)
     # Create the bucket if it does not exist.
    bucket = client.create_bucket(bucket, location="us-central1")
    # Enable uniform bucket-level access
    bucket.iam_configuration.uniform_bucket_level_access_enabled = True
    bucket.patch()

    # Grant public read access to the objects in the bucket
    policy = bucket.get_iam_policy(requested_policy_version=3)
    policy.bindings.append({
        "role": "roles/storage.objectViewer",
        "members": {"allUsers"}
    })
    bucket.set_iam_policy(policy)

if __name__ == '__main__':
    url = sys.argv[1]
    id = sys.argv[2]
    scrape_page(url, id)
    create_bucket(id)