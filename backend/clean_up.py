import os
from google.cloud import storage

def clean_up():
    os.environ['GO a GOOGLE_APPLICATION_CREDENTIALS'] = './config/credentials.json'

    # Create a Cloud Storage client.
    client = storage.Client()

    # Get the bucket.
    bucket = client.bucket("translated-images")

    # List all objects in the bucket, set cache-control, and delete them.
    blobs = bucket.list_blobs()
    for blob in blobs:
        # Delete the blob
        blob.delete()

    print("Finished")

if __name__ == '__main__':
    clean_up()
