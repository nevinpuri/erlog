FROM python:3.11-slim

EXPOSE 8000

RUN pip install fastapi ujson structlog chdb uvicorn luqum async_tail

COPY ./app/ ./app/

CMD ["python3", "-m", "uvicorn", "app.main:app"]
