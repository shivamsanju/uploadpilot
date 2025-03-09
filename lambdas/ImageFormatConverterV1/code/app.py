import boto3
import io
from PIL import Image
from event_helper import enforce_input_output_format, InputEvent, OutputFormat


ALLOWED_EXTENSIONS = ["png", "jpg", "jpeg", "webp", "bmp", "gif", "tiff"]


def lambda_handler(event, context):
    return image_convertor(event, context)


@enforce_input_output_format
def image_convertor(event: InputEvent, context) -> OutputFormat:
    try:

        ext = event.input_filename.split(".")[-1]
        if ext not in ALLOWED_EXTENSIONS:
            raise ValueError(f"Unsupported file extension. Only {ALLOWED_EXTENSIONS} are supported.")
    
        requested_format = event.input_args.get("format", "")
        if requested_format not in ALLOWED_EXTENSIONS:
            raise ValueError(f"Unsupported file extension. Only {ALLOWED_EXTENSIONS} are supported.")
    
        s3_client = boto3.client("s3")

        response = s3_client.get_object(Bucket=event.input_bucket, Key=event.input_key)
        image_data = response["Body"].read()

        image = Image.open(io.BytesIO(image_data))
        output_buffer = io.BytesIO()
        image.save(output_buffer, format=requested_format.upper())

        new_filename = event.input_filename.rsplit(".", 1)[0] + "." + requested_format
        new_key = event.output_key_prefix + new_filename

        # Upload the PNG image back to S3
        s3_client.put_object(
            Bucket=event.output_bucket,
            Key=new_key,
            Body=output_buffer.getvalue(),
            ContentType=f"image/{requested_format}"
        )

        return OutputFormat(
            status_code=200,
            message="Success",
            bucket=event.output_bucket,
            output_filename=new_filename,
            output_key=new_key,
            output_content_type=f"image/{requested_format}"
        )

    except Exception as e:
        return OutputFormat(
            status_code=500,
            message="Error",
            error=str(e)
        )