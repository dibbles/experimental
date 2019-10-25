# Architecture Information

1. [Changing The Polling Duration](#webhook-runtime-architecture)
2. [Overriding The Status Message](#webhook-creation-architecture)

## Webhook Runtime Architecture

The following diagram shows what occurs at runtime when a webhooks are triggered.

![Architecture Diagram](./images/architecture.png?raw=true "Diagram showing overall runtime architecture of the webhooks extension")


# tekton-validate-github-event

For https://github.com/tektoncd/experimental/issues/245

Checks the following

1) Valid X-Hub signature (secret token used in validation service matches the secret token used on the webhook)
2) Repository URL matches the input URL parameter - so we only activate Triggers for selected repositories
3) Eventually - that it's only for a push or pull request event


## Webhook Creation Architecture