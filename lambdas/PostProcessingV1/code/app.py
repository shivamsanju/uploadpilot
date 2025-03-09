import boto3
import io
import zipfile
from event_helper import OutputFormat



def lambda_handler(event, context):
    s3 = boto3.client('s3')
    
    try :
        bucket = event["workspace_id"]
        file_keys, delete_callback = get_processed_folder_files(event)

        output_key = f'{event["upload_id"]}/artifacts/{event["processor_id"]}/{event["run_id"]}.zip'
        
        zip_buffer = io.BytesIO()
        with zipfile.ZipFile(zip_buffer, 'w', zipfile.ZIP_DEFLATED) as zip_file:
            for file_key in file_keys:
                file_obj = s3.get_object(Bucket=bucket, Key=file_key)
                file_data = file_obj['Body'].read()
                zip_file.writestr(file_key.split('/')[-2] + '/' + file_key.split('/')[-1], file_data)
        
        zip_buffer.seek(0)
        
        # Upload zip to S3
        s3.put_object(Bucket=bucket, Key=output_key, Body=zip_buffer.getvalue())

        # TODO: Handle if error occurs here
        delete_callback()
        delete_staging_folder(event)
        
        return OutputFormat(
            status_code=200,
            message="Successfully zipped artifacts",
            bucket=bucket,
            output_key=output_key,
            output_filename=f'{event["run_id"]}.zip',
            output_content_type="application/zip",
            error="",
        ).to_dict()
    
    except Exception as e:
        return OutputFormat(
            status_code=500,
            message="The activity did not return a valid response.",
            error=str(e)
        ).to_dict()
    


def get_processed_folder_files(event):
    s3 = boto3.client('s3')
    source_prefix = f'{event["upload_id"]}/processed/{event["processor_id"]}/{event["run_id"]}/'
    
    response = s3.list_objects_v2(Bucket=event["workspace_id"], Prefix=source_prefix)
    
    file_keys = []
    delete_objects = []
    contents = response.get('Contents', [])
    for obj in contents:
        file_key = obj['Key']
        delete_objects.append({'Key': file_key})
        if file_key.endswith('/'):
            continue
        file_keys.append(file_key)
    
    def delete_folder_callback():
        if delete_objects:
            s3.delete_objects(Bucket=event["workspace_id"], Delete={'Objects': delete_objects})
    
    return file_keys, delete_folder_callback


def delete_staging_folder(event):
    s3 = boto3.client('s3')
    source_prefix = f'{event["upload_id"]}/staging/{event["processor_id"]}/{event["run_id"]}/'
    
    response = s3.list_objects_v2(Bucket=event["workspace_id"], Prefix=source_prefix)
    
    contents = response.get('Contents', [])
    
    if contents:
        delete_objects = [{'Key': obj['Key']} for obj in contents]
        s3.delete_objects(Bucket=event["workspace_id"], Delete={'Objects': delete_objects})