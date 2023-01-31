# Tiltfile
load("ext://helm_remote", "helm_remote")
load("ext://secret", "secret_from_dict")

docker_build('transaction-card-validator', '.')
docker_build("e2e-tests", ".", dockerfile="e2e/e2e.Dockerfile")

# Minio
helm_remote(
    'minio',
    repo_name='bitnami',
    repo_url='https://charts.bitnami.com/bitnami',
    set=[
        "auth.rootUser=minio",
        "auth.rootPassword=minio123",
        "defaultBuckets=transactions-way4-quarantine-dev transactions-solar-quarantine-dev transactions-test-quarantine-dev reconciliation-solar-quarantine-dev",
        "persistence.enabled=false",
        "persistence.storageClass=standard"
    ]
)

k8s_resource(workload='minio', port_forwards='9001', labels=["minio"])

# Kafka
helm_remote(
    "kafka",
    repo_name="bitnami",
    repo_url="https://charts.bitnami.com/bitnami",
    release_name="kafka",
    set=[
        # "auth.clientProtocol=sasl",
        # "auth.interBrokerProtocol=sasl",
        # "auth.sasl.mechanisms=plain",
        # "auth.sasl.interBrokerMechanism=plain",
        # "auth.sasl.jaas.clientUsers[0]=username",
        # "auth.sasl.jaas.clientPasswords[0]=password",
        # "auth.sasl.jaas.interBrokerUser=admin",
        # "auth.sasl.jaas.interBrokerPassword=admin",
        "listeners[0]=INTERNAL://:9092",
        "listeners[1]=CLIENT://localhost:9093",
        "advertisedListeners[0]=INTERNAL://:9092",
        "advertisedListeners[1]=CLIENT://localhost:9093",
        "interBrokerListenerName=INTERNAL"
  ]
)


k8s_resource("kafka", port_forwards='9093', labels=["kafka"])
k8s_resource("kafka-zookeeper", labels=["kafka"])

# Simple schema-registry
k8s_yaml(
    helm(
        'charts/local',
        name='saltdata',
        values=['charts/local/values.yaml']
    )
)
k8s_resource(workload='saltdata-cp-schema-registry', port_forwards='8081', labels=["kafka"], resource_deps=["kafka"])

# Transaction-card-validator
k8s_yaml(
    helm(
        "charts/transaction-card-validator",
        name="transaction-card-validator",
        values=["charts/local/values/transaction.yaml"],
        set= [
            "validator.image.repository=transaction-card-validator",
            "validator.image.tag=latest"
        ]
    )
)

k8s_resource("transaction-card-validator-test-ingress", labels=["validator"], resource_deps=["saltdata-cp-schema-registry"])

k8s_yaml(
    secret_from_dict(
        "transaction-card-validator-msk-eventstreaming",
        inputs={
            "endpoint": "kafka-0.kafka-headless.default.svc.cluster.local:9092",
            "username": "",
            "password": "",
        },
    )
)

k8s_yaml('./e2e/e2e.yaml' )
k8s_resource('e2e-tests',
   resource_deps=['kafka'],
   trigger_mode=TRIGGER_MODE_MANUAL,
   auto_init=False,
   labels=["tests"],
)

# Run tests button
load('ext://uibutton', 'cmd_button', 'location', 'text_input')
cmd_button(name='run-command',
           resource='transaction-card-validator-test-ingress',
           inputs=[
             text_input('COMMAND', label="command", default="make unit-tests integration-tests acceptance-tests"),
           ],
           argv=['sh', '-c', '$COMMAND'],
           text='Run tests',
           icon_name='science')