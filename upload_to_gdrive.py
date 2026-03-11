import argparse
import os
import sys

try:
    from googleapiclient.discovery import build
    from googleapiclient.http import MediaFileUpload
    from google.oauth2 import service_account
except ImportError:
    print("Please install required packages: pip install google-api-python-client google-auth-httplib2 google-auth-oauthlib")
    sys.exit(1)

SCOPES = ['https://www.googleapis.com/auth/drive.file', 'https://www.googleapis.com/auth/drive']

def authenticate(service_account_file):
    creds = service_account.Credentials.from_service_account_file(
        service_account_file, scopes=SCOPES)
    return build('drive', 'v3', credentials=creds)

def find_file_in_folder(service, file_name, folder_id):
    query = f"name='{file_name}' and trashed=false"
    if folder_id:
        query += f" and '{folder_id}' in parents"

    results = service.files().list(q=query, spaces='drive', fields='nextPageToken, files(id, name)').execute()
    items = results.get('files', [])
    if not items:
        return None
    return items[0]['id']

def upload_or_update_file(service, file_path, folder_id=None):
    file_name = os.path.basename(file_path)
    file_metadata = {'name': file_name}

    existing_file_id = find_file_in_folder(service, file_name, folder_id)

    media = MediaFileUpload(file_path, mimetype='application/json', resumable=True)

    if existing_file_id:
        print(f"Updating existing file {file_name} (ID: {existing_file_id})...")
        file = service.files().update(fileId=existing_file_id, media_body=media).execute()
        print(f"File ID: {file.get('id')} successfully updated.")
    else:
        print(f"Creating new file {file_name}...")
        if folder_id:
            file_metadata['parents'] = [folder_id]
        file = service.files().create(body=file_metadata, media_body=media, fields='id').execute()
        print(f"File ID: {file.get('id')} successfully uploaded.")

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Upload or update a JSON file to Google Drive using a Service Account.")
    parser.add_argument("file", help="Path to the JSON file to upload")
    parser.add_argument("--credentials", default="credentials.json", help="Path to your service account credentials JSON file (default: credentials.json)")
    parser.add_argument("--folder", help="Optional Google Drive folder ID to upload into")

    args = parser.parse_args()

    if not os.path.exists(args.file):
        print(f"Error: File not found: {args.file}")
        sys.exit(1)

    if not os.path.exists(args.credentials):
        print(f"Error: Credentials file not found: {args.credentials}")
        print("Please obtain a service account JSON file from Google Cloud Console and save it here.")
        sys.exit(1)

    service = authenticate(args.credentials)
    upload_or_update_file(service, args.file, args.folder)
