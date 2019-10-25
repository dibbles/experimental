# Architecture Information

## Contents

1. [End User Overview](#end-user-overview)
2. [Webhook Runtime Architecture](#webhook-runtime-architecture)
2. [Webhook Creation Architecture](#webhook-creation-architecture)

 
## End User Overview

![User Setup Diagram](./images/setup.png?raw=true "Diagram showing initial user setup")

The diagram above shows the initial configuration required prior to using the webhooks extension.  The process consists of:

1) Installing the necessary parts of Tekton

2) Creating a secret (if using ingress), that holds the SSL certificate securing the endpoint.  Created in the same namespace as the Tekton installation.

3) Installing the triggertemplate for your pipeline into the same namespace as the Tekton installation. *Your triggertemplate must currently be named <pipeline-name>-template*.

4) Installing the triggerbindings for your pipeline into the same namespace as the Tekton installation. *You need to install two triggerbindings per pipeline*, one for pull request events and one for push events.  It is believed that this is likely to be changed in a future release. *Your triggerbindings must currently be named <pipeline-name>-push-binding and <pipeline-name>-pullrequest-binding*.

5) Creating secrets for accessing GitHub and Docker and patching them onto the service account under which you want your pipelinerun to execute.  These secrets need creating in the *target namespace* where you want your pipeline to run.  Note that some of this process can be completed using the Tekton Dashboard.

6) Installing the pipeline into the *target namespace* where you want the pipeline to run.

 
## Webhook Runtime Architecture

![Architecture Diagram](./images/architecture.png?raw=true "Diagram showing overall runtime architecture of the webhooks extension")

The diagram above shows what occurs at runtime when webhooks are triggered.  The process consists of:

1) All webhooks communicate with a single ingress/route as an access point to the cluster.

2) The ingress/route is backed by the eventlistener service.  The eventlistener iterates over all of the triggers defined in the eventlistener custom resource (labeled 3) and sends a request for each trigger to the interceptor service (labeled 4).

3) The interceptor section for each trigger within the eventlistener contains conditions under which that trigger should operate.  For example, the event is from git repository X and the event is a pull_request.

4) The interceptor service's response to each request determines whether or not the trigger is valid for the incoming webhook event.  The interceptor checks:

 - Valid X-Hub signature - secret token defined at webhook creation matches the secret token on the incoming webhook.
 - Repository URL matches - so we only activate a trigger for a selected repository.
 - Webhook event matches - so we only activate a trigger for a selected event type, a push or pull request event.

5) The Tekton Triggers code creates the necessary pipelineresources, pipelineruns etc... as defined in the triggertemplate - substituting parameters as defined in the triggerbinding or from the parameters set on the trigger in the eventlistener.

In the case that the event type is a pull request, a monitor taskrun will be created to monitor the pipelineruns and report status onto the pull request in GitHub.

 
## Webhook Creation Architecture