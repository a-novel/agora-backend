# Agora des Écrivains

Private repository for Agora backend application.

# Run

```shell
make run
```

> You'll need some signature keys in order to authenticate. To generate some locally, run this command (with the local
> server started):
>  ```shell
>  make rotate-keys
>  ```

# Test

```shell
# Optional, if your local environment is not up.
# docker compose up -d

make test
# Or `make race` to run in race mode.
```

> `make msan` command only works with some Linux instances, so you cannot run it on your mac.
> Check the CI to ensure it passes properly.

# Monitor

You can connect to your local database by running the following command:

```shell
make db
```

# Deployment

Agora runs on Google Cloud Platform. You will need to be added to the project to be able to monitor deployment.

Below is the current stack of resources currently used by the project:

- [Cloud Run](https://console.cloud.google.com/run?project=agoradesecrivains)
  for deploying code.
- [Cloud Storage](https://console.cloud.google.com/storage/browser?project=agoradesecrivains)
  to emulate a local filesystem.
  - `backend-token-keys`: stores signature keys for JWT tokens.
- [VPC Networks](https://console.cloud.google.com/networking/networks/list?project=agoradesecrivains)
  for securing connections between services.
  - `backend`: only allows traffic from the backend application.
- [Cloud SQL](https://console.cloud.google.com/sql/instances?project=agoradesecrivains)
  for storing data.
  - `agora-postgres`: main instance used by this service.
- [Secret Manager](https://console.cloud.google.com/security/secret-manager?project=agoradesecrivains)
  for storing sensitive data.
  - `sendgrid-api-key`: stores the Sendgrid API key.
  - `agora-postgres-dsn`: dsn string to access the main database.
- [Load Balancer](https://console.cloud.google.com/net-services/loadbalancing/list/loadBalancers?project=agoradesecrivains)
  redirects traffic from the domain to the backend application.
- [Cloud Scheduler](https://console.cloud.google.com/cloudscheduler?project=agoradesecrivains)
  to schedule backend jobs.

You can look at 
[production logs here](https://console.cloud.google.com/logs/query;query=jsonPayload.app%3D%22agora-backend%22?project=agoradesecrivains).
