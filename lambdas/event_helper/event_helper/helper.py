import functools
from dataclasses import dataclass,field


def enforce_input_output_format(func):
    return enforce_output_format(enforce_input_event(func))


@dataclass
class InputEvent:
    activity_key: str = ""
    input_bucket: str = ""
    input_key: str = ""
    input_filename: str = ""
    input_content_type: str = ""
    input_args: dict = field(default_factory=dict)  # Correct way to define a mutable default
    output_bucket: str = ""
    output_key_prefix: str = ""

def enforce_input_event(func):
    @functools.wraps(func)
    def wrapper(event, *args, **kwargs):
        try:
            check_required_keys(event)
            input_activity_exists, input_activity_key = check_input_activity_key(event)
            
            input_event = InputEvent(
                activity_key=event["current_activity_key"],
            )

            if input_activity_exists:
                input_event.input_bucket = event[f"{input_activity_key}.bucket"]
                input_event.input_key = event[f"{input_activity_key}.output_key"]
                input_event.input_filename = event[f"{input_activity_key}.output_filename"]
                input_event.input_content_type = event[f"{input_activity_key}.output_content_type"]
            else:
                input_event.input_bucket = event["workspace_id"]
                input_event.input_key = f"{event['upload_id']}/raw/{event['file_name']}"
                input_event.input_filename = event["file_name"]
                input_event.input_content_type = event["content_type"]


            input_event.output_bucket = event["workspace_id"]
            input_event.output_key_prefix = get_output_prefix(event)
            input_event.input_args = get_input_args(event)
            
            return func(input_event, *args, **kwargs)

        except Exception as e:
            return OutputFormat(
                status_code=500,
                message="Error",
                error="The activity does not have a valid input: " + str(e)
            )
        
    return wrapper


def check_required_keys(event):
    required_keys = {"workspace_id", "upload_id", "processor_id", "run_id", "file_name", "content_type", "current_activity_key"}
    if not all(key in event for key in required_keys):
        raise ValueError("Missing required keys in event.")


def check_input_activity_key(event) -> bool:
    if f"{event['current_activity_key']}.input" in event and event[f"{event['current_activity_key']}.input"] != "":
        input_activity_key = event[f"{event['current_activity_key']}.input"]

        if f"{input_activity_key}.bucket" not in event:
            raise ValueError(f"Missing required keys in event : {f'{input_activity_key}.bucket'}")

        if f"{input_activity_key}.output_key" not in event:
            raise ValueError(f"Missing required keys in event : {f'{input_activity_key}.output_key'}")
        
        if f"{input_activity_key}.output_filename" not in event:
            raise ValueError(f"Missing required keys in event : {f'{input_activity_key}.output_filename'}")
        
        if f"{input_activity_key}.output_content_type" not in event:
            raise ValueError(f"Missing required keys in event : {f'{input_activity_key}.output_content_type'}")

        return True, input_activity_key

    return False, ""


def get_output_prefix(event):
    output_folder = "staging"
    if f"{event['current_activity_key']}.save_output" in event and event[f"{event['current_activity_key']}.save_output"] == True:
        output_folder = "processed"

    return f'{event["upload_id"]}/{output_folder}/{event["processor_id"]}/{event["run_id"]}/{event["current_activity_key"]}/'

def get_input_args(event):
    input_args = {}
    for key, value in event.items():
        if key.split(".")[0] == event["current_activity_key"]:
            input_args[key.split(".")[1]] = value
    return input_args

################################ Output Format ################################

@dataclass
class OutputFormat:
    status_code: int = 200
    message: str = "Success"
    bucket: str = ""
    output_key: str = ""
    output_filename: str = ""
    output_content_type: str = ""
    error: str = ""

    def to_dict(self) -> dict:
        return {
            "status_code": self.status_code,
            "message": self.message,
            "bucket": self.bucket,
            "output_key": self.output_key,
            "output_filename": self.output_filename,
            "output_content_type": self.output_content_type,
            "error": self.error,
        }

def enforce_output_format(func):
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        try:
            result = func(*args, **kwargs)
            if not isinstance(result, OutputFormat) or not result.status_code:
                raise ValueError("Function must return an instance of OutputFormat with status_code set.")
            return result.to_dict()
        except Exception as e:
            return OutputFormat(
                status_code=500,
                message="The activity did not return a valid response.",
                error=str(e)
            ).to_dict()
    
    return wrapper



