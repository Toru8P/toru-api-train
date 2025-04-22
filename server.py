# server.py
from fastapi import FastAPI
import docker
import uuid

app = FastAPI()
client = docker.from_env()

@app.post("/train/lora")
def train_lora():
    job_id = str(uuid.uuid4())
    container = client.containers.run(
        image="lora-trainer:latest",
        name=f"lora-job-{job_id}",
        detach=True,
        volumes={
            "/shared/logs": {"bind": "/app/logs", "mode": "rw"}
        }
    )
    return {"job_id": job_id, "container_id": container.id}

@app.get("/logs/lora/{job_id}")
def get_logs(job_id: str):
    try:
        container = client.containers.get(f"lora-job-{job_id}")
        logs = container.logs(tail=100).decode()
        return {"logs": logs}
    except docker.errors.NotFound:
        return {"error": "Job not found"}

@app.get("/status/lora/{job_id}")
def check_status(job_id: str):
    try:
        container = client.containers.get(f"lora-job-{job_id}")
        return {"status": container.status}
    except docker.errors.NotFound:
        return {"error": "Job not found"}

@app.post("/cancel/lora/{job_id}")
def cancel_job(job_id: str):
    try:
        container = client.containers.get(f"lora-job-{job_id}")
        container.stop()
        return {"status": "cancelled"}
    except docker.errors.NotFound:
        return {"error": "Job not found"}
