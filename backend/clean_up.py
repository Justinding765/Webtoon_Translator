import os
from google.cloud import storage
import sys

def clean_up(sessionID):
    os.environ['GOOGLE_APPLICATION_CREDENTIALS'] = './config/credentials.json'

    # Create a Cloud Storage client.
    client = storage.Client()

    # Get the bucket with the name sessionID.
    bucket = client.bucket(sessionID)

    # Check if the bucket exists.
    if bucket.exists():
        # List all objects in the bucket and delete them.
        blobs = bucket.list_blobs()
        for blob in blobs:
            # Delete the blob
            blob.delete()
        bucket.delete()

    script_dir = os.path.dirname(os.path.abspath(__file__))
    output_folder = os.path.join(script_dir, 'output')
    html_file = os.path.join(output_folder, f'{sessionID}.html')

    if os.path.exists(html_file):
        os.remove(html_file) 
    print("finished cleaning")

if __name__ == '__main__':

    sessionID = sys.argv[1]
    clean_up(sessionID)
