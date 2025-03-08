import json
import boto3
import logging

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def lambda_handler(event, context):

    verify_supported_input(event)
    
    bucket_name, object_key, filename = get_file_key(event)

    try:
        s3_client = boto3.client("s3")
        textract_client = boto3.client('textract')

        response = textract_client.detect_document_text(
            Document={
                'S3Object': {
                    'Bucket': bucket_name,
                    'Name': object_key
                }
            }
        )
                # Extract detected text
        detected_text = []
        for item in response.get('Blocks', []):
            if item['BlockType'] == 'LINE':
                detected_text.append(item['Text'])
        
        # Join detected text into a single string
        text_output = '\n'.join(detected_text)


        new_filename = filename.rsplit(".", 1)[0] + ".txt"
        new_key = get_output_key(event, new_filename)

        # Upload the PNG image back to S3
        s3_client.put_object(
            Bucket=bucket_name,
            Key=new_key,
            Body=text_output,
            ContentType="image/png"
        )

        return {
            "success": True,

        }

    except Exception as e:
        return {
            "success": False,
            "error": str(e),
        }
    



ALLOWED_EXTENSIONS = ["png", "jpg", "jpeg", "webp", "bmp", "gif", "tiff", "pdf", "doc", "docx"]


def verify_supported_input(event):
    filename = event["file_name"]
    ext = filename.split(".")[-1]

    if ext not in ALLOWED_EXTENSIONS:
        raise ValueError(f"Unsupported file extension. Only {ALLOWED_EXTENSIONS} images are supported.")
    

def get_file_key(event):
    bucket = event["workspace_id"]
    filename = event["file_name"]
    upload_id = event["upload_id"]

    if bucket == "" or filename == "" or upload_id == "":
        raise ValueError("Invalid event. Missing required fields.")
    
    return bucket, f'{upload_id}/raw/{filename}', filename
 

def get_output_key(event, new_file_name):
    upload_id = event["upload_id"]
    processor_id = event["processor_id"]
    run_id = event["run_id"]
    activity_key = event["activity_key"]

    if new_file_name == "" or upload_id == "" or processor_id == "" or run_id == "":
        raise ValueError("Invalid event. Missing required fields.")
    
    return f'{upload_id}/processed/{processor_id}/{run_id}/{activity_key}/{new_file_name}'