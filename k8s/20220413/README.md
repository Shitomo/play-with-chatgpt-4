# Prompt

Write k8s manifest which define below
1. Deployment which runs a sample api server image.
2. Service which exposes the Deployment
3. Deployment which runs a sample api server image which access 1. service
4. Open Telemetry Collector which collect metrics from 1. and 3. Deployment


# Answer

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
  labels:
    app: api-server
spec:
    replicas: 1
    selector:
        matchLabels:
        app: api-server
    template:
        metadata:
        labels:
            app: api-server
        spec:
        containers:
        - name: api-server
            image: quay.io/kenfdev/api-server:latest
            ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: api-server
spec:
    selector:
        app: api-server
    ports:
    - protocol: TCP
        port: 80
        targetPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-client
  labels:
    app: api-client
spec:
    replicas: 1
    selector:
        matchLabels:
        app: api-client
    template:
        metadata:
        labels:
            app: api-client
        spec:
        containers:
        - name: api-client
            image: quay.io/kenfdev/api-client:latest
            ports:
            - containerPort: 8080
            env:
            - name: API_SERVER_URL
                value: http://api-server
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otel-collector
  labels:
    app: otel-collector
spec:
    replicas: 1
    selector:
        matchLabels:
        app: otel-collector
    template:
        metadata:
        labels:
            app: otel-collector
        spec:
        containers:
        - name: otel-collector
            image: otel/opentelemetry-collector-contrib:0.38.0
            ports:
            - containerPort: 55680
            - containerPort: 55681
            - containerPort: 8888
            - containerPort: 8889
            - containerPort: 9411
            - containerPort: 14250
            - containerPort: 14268
    ---
    apiVersion: v1
    kind: Service
    metadata:
    name: otel-collector
    spec:
    selector:
        app: otel-collector
    ports:
    - protocol: TCP
        port: 55680
        targetPort: 55680
    - protocol: TCP
        port: 55681
        targetPort: 55681
    - protocol: TCP
        port: 8888
        targetPort: 8888
    - protocol: TCP
        port: 8889
        targetPort: 8889
    - protocol: TCP 
        port: 9411
        targetPort: 9411
    - protocol: TCP 
        port: 14250
        targetPort: 14250
    - protocol: TCP
        port: 14268
        targetPort: 14268

    
    ```
