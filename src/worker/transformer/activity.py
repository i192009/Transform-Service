from temporalio import activity

@activity.defn(name="DownloadResource")
async def DownloadResource(token: str, downloadUrl: str) -> bool:
    # download one resource
    return True

@activity.defn(name="ExecuteCommandLine")
async def ExecuteCommandLine(job: object) -> bool:
    # Execute command here
    return True

@activity.defn(name="UploadResource")
async def UploadResource(token: str, localPath: str, remoteUrl: str) -> bool:
    # upload one resource
    return True