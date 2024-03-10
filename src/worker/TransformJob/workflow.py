# workflows.py

from temporal.workflow import workflow_method, Workflow

class JobWorkflow:
    @workflow_method
    async def handle_job(self, job_id):
        # Call activities in the required sequence
        await download_resources(job_id)
        await process_job(job_id)
        await upload_results(job_id)

    @workflow_method
    async def cancel_job_workflow(self, job_id):
        # Workflow to handle job cancellation
        await cancel_job(job_id)
