from datetime import timedelta
from temporalio import workflow
from temporalio.common import RetryPolicy
import json

@workflow.defn
class FileConvert:
    @workflow.run
    async def run(self, job: str) -> str:
        job = json.loads(job)
        for file in job.Files:
            if file.Transform:
                await workflow.execute_activity(
                    "DownloadResource",
                    [job.Token, file.DownloadKey],
                    schedule_to_close_timeout=timedelta(seconds=5),
                    # Retry Policy
                    retry_policy=RetryPolicy(
                        backoff_coefficient=2.0,
                        maximum_attempts=5,
                        initial_interval=timedelta(seconds=1),
                        maximum_interval=timedelta(seconds=2),
                        # non_retryable_error_types=["ValueError"],
                    ),
                )
                # transform file
                # upload file
                # update file.UploadKey
        
        await workflow.execute_activity(
            "ExecuteCommandLine",
            job,
            schedule_to_close_timeout=timedelta(seconds=5),
            # Retry Policy
            retry_policy=RetryPolicy(
                backoff_coefficient=2.0,
                maximum_attempts=5,
                initial_interval=timedelta(seconds=1),
                maximum_interval=timedelta(seconds=2),
                # non_retryable_error_types=["ValueError"],
            ),
        )

        for file in job.UploadFiles:
            await workflow.execute_activity(
                "UploadResource",
                file.UploadKey,
                schedule_to_close_timeout=timedelta(seconds=5),
                # Retry Policy
                retry_policy=RetryPolicy(
                    backoff_coefficient=2.0,
                    maximum_attempts=5,
                    initial_interval=timedelta(seconds=1),
                    maximum_interval=timedelta(seconds=2),
                    # non_retryable_error_types=["ValueError"],
                ),
            )
            # update file.UploadKey
