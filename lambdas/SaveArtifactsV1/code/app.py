import boto3
import io
import zipfile

def lambda_handler(event, context):
    s3 = boto3.client('s3')
    try :
        workspace_id = event["workspace_id"]
        upload_id = event["upload_id"]
        processor_id = event["processor_id"]
        run_id = event["run_id"]

        if upload_id == "" or processor_id == "" or run_id == "":
            raise ValueError("Invalid event. Missing required fields.")

        
        source_prefix = f'{upload_id}/processed/{processor_id}/{run_id}/'
        destination_key = f'{upload_id}/processed-zip/{processor_id}/{run_id}.zip'
        
        response = s3.list_objects_v2(Bucket=workspace_id, Prefix=source_prefix)
        
        if 'Contents' not in response:
            return {'success': False, 'error': 'No files found to zip'}
        
        zip_buffer = io.BytesIO()
        
        with zipfile.ZipFile(zip_buffer, 'w', zipfile.ZIP_DEFLATED) as zip_file:
            for obj in response['Contents']:
                file_key = obj['Key']
                if file_key.endswith('/'):
                    continue  # Skip directories
                
                file_obj = s3.get_object(Bucket=workspace_id, Key=file_key)
                file_data = file_obj['Body'].read()
                
                zip_file.writestr(file_key[len(source_prefix):], file_data)  # Trim prefix for better structure
        
        zip_buffer.seek(0)
        
        # Upload zip to S3
        s3.put_object(Bucket=workspace_id, Key=destination_key, Body=zip_buffer.getvalue())
        
        return {'success': True, 'message': 'Files zipped and uploaded successfully'}
    except Exception as e:
        return {'success': False, 'error': str(e)}
    