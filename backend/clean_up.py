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
    # Create the bucket if it does not exist.
    bucket = client.create_bucket(bucket, location="us-central1")
    print(f"Bucket {sessionID} created.")

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

    print(sessionID)

if __name__ == '__main__':
    sessionID = sys.argv[1]
    clean_up(sessionID)
