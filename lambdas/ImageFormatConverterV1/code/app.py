import boto3
import io
from PIL import Image

s3_client = boto3.client("s3")

def lambda_handler(event, context):

    verify_supported_input(event)
    
    bucket_name, object_key, filename = get_file_key(event)

    requested_format = event.get(f'{event.get("activity_key", "")}_format', "")

    if requested_format not in ALLOWED_EXTENSIONS:
        raise ValueError(f"Unsupported file extension. Only {ALLOWED_EXTENSIONS} are supported.")

    try:
        response = s3_client.get_object(Bucket=bucket_name, Key=object_key)
        image_data = response["Body"].read()

        image = Image.open(io.BytesIO(image_data))
        output_buffer = io.BytesIO()
        image.save(output_buffer, format=requested_format.upper())

        new_filename = filename.rsplit(".", 1)[0] + "." + requested_format
        new_key = get_output_key(event, new_filename)

        # Upload the PNG image back to S3
        s3_client.put_object(
            Bucket=bucket_name,
            Key=new_key,
            Body=output_buffer.getvalue(),
            ContentType=f"image/{requested_format}"
        )

        return {
            "success": True,
            "file_name": new_filename,
            "content_type": f"image/{requested_format}",
            "output_key": new_key,
            "message": "Conversion successful"
        }

    except Exception as e:
        return {
            "success": False,
            "error": str(e),
        }
    



ALLOWED_EXTENSIONS = ["png", "jpg", "jpeg", "webp", "bmp", "gif", "tiff"]

ALLOWED_CONTENT_TYPES = [
    "image/png",
    "image/jpeg",
    "image/jpeg",
    "image/webp",
    "image/bmp",
    "image/gif",
    "image/tiff"
]


def verify_supported_input(event):
    content_type = event["content_type"]
    filename = event["file_name"]
    ext = filename.split(".")[-1]

    if ext not in ALLOWED_EXTENSIONS:
        raise ValueError(f"Unsupported file extension. Only {ALLOWED_EXTENSIONS} are supported.")
    
    if content_type not in ALLOWED_CONTENT_TYPES:
        raise ValueError(f"Unsupported content type. Only {ALLOWED_CONTENT_TYPES} are supported.")


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