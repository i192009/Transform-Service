# worker.py

from temporal.worker import WorkerFactory
import activities
import workflows

# Configure the Temporal worker factory
factory = WorkerFactory("localhost:7233", namespace="default")
worker = factory.new_worker("default")
worker.register_activities_implementation(activities, "JobActivities")
worker.register_workflow_implementation_type(workflows.JobWorkflow)

factory.start()
