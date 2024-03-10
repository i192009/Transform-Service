import asyncio
from temporalio.client import Client
from temporalio.worker import Worker
from transformer.workflow import ConvertFile
from transformer.activity import *
# ...
# ...
async def main():
    client = await Client.connect("localhost:7233")
    worker = Worker(
        client,
        task_queue="your-task-queue",
        workflows=[ConvertFile],
        activities=[DownloadResource, ExecuteCommandLine, UploadResource],
    )
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())