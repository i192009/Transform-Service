from concurrent import futures
import grpc
import TransformService2_pb2
import TransformService2_pb2_grpc
from temporal.workflow import WorkflowClient

class TransformV2Service(transform_v2_pb2_grpc.TransformV2ServiceServicer):

    def __init__(self):
        self.client = WorkflowClient.new_client(namespace="default")

    def CreateJob(self, request, context):
        # Start the job workflow
        workflow = self.client.new_workflow_stub(workflows.JobWorkflow)
        workflow.handle_job(request.jobId)
        return transform_v2_pb2.CreateJobResponse(jobId=request.jobId)

    def CancelJob(self, request, context):
        # Start the cancel job workflow
        workflow = self.client.new_workflow_stub(workflows.JobWorkflow)
        workflow.cancel_job_workflow(request.jobId)
        return transform_v2_pb2.CancelJobResponse(status="Cancellation Initiated")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    transform_v2_pb2_grpc.add_TransformV2ServiceServicer_to_server(TransformV2Service(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
