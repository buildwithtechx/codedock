---
title: AI Workloads
description: Run model and vector services as explicitly configured containers.
---

AI workloads use the same image, variable, volume and networking settings as other application services. They are recipes, not an automatic AI provisioning system.

## Ollama

Create an image-based service using a versioned Ollama image. Configure its listening port and a persistent volume for downloaded models. Allow for model download time and host memory requirements before calling the service from another application.

Keep the model endpoint private unless you add appropriate access controls. GPU use requires a compatible host, drivers and Docker runtime configuration; Codedock does not automatically provision GPU hardware or configure those dependencies.

## Qdrant

Use a versioned Qdrant image, persist its storage directory and configure authentication according to the image's requirements. Keep its HTTP and gRPC endpoints restricted to the applications that need them. Back up vector data using a method compatible with your Qdrant version.

## PostgreSQL with pgvector

Use a PostgreSQL image containing the extension, then enable it with your migration tooling. A standard PostgreSQL service does not guarantee pgvector is installed. Confirm image compatibility before restoring an existing database.

## Codedock AI settings

AI provider settings used by Codedock for diagnostic assistance are separate from the AI containers you deploy. Configure the provider URL, model and credentials in settings only if you want that integration. Installing a model service does not automatically connect it as Codedock's AI provider.

See [templates and recipes](/deployments/templates/) and [build strategies](/deployments/build-strategies/).
