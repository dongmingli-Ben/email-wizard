# Use the official Python image as the base image
FROM python:latest

WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt

WORKDIR /mnt
RUN rm -rf /app