# Install on Red Hat OpenShift

Assuming you've completed the [prereq installation and setup](./InstallPrereqs.md),

1. If you plan to use `buildah` in your Pipelines, you will need to set an additional permission on any service account that will be used to run a pipeline by using the following command:

      ```
      oc adm policy add-scc-to-user privileged -z [service_account_name] -n [namespace]
      ```

2. Install the webhooks extension:

      - Install the [release build](./InstallReleaseBuild.md)
      - Install the [nightly build](./InstallNightlyBuild.md)

3. Check you can access the Webhooks Extension through the Dashboard UI that you should already have a Route for, for example at http://tekton-dashboard.${openshift_master_default_subdomain}/#/extensions/webhooks-extension.

    ![Create webhook page in dashboard](./images/createWebhook.png?raw=true "Create webhook page in dashboard")

4. Begin creating webhooks


## Notes:

This has been tested with the following scc (from `oc get scc`):

```
NAME               PRIV      CAPS      SELINUX     RUNASUSER          FSGROUP     SUPGROUP    PRIORITY   READONLYROOTFS   VOLUMES
anyuid             false     []        MustRunAs   RunAsAny           RunAsAny    RunAsAny    10         false            [configMap downwardAPI emptyDir persistentVolumeClaim projected secret]
hostaccess         false     []        MustRunAs   MustRunAsRange     MustRunAs   RunAsAny    <none>     false            [configMap downwardAPI emptyDir hostPath persistentVolumeClaim projected secret]
hostmount-anyuid   false     []        MustRunAs   RunAsAny           RunAsAny    RunAsAny    <none>     false            [configMap downwardAPI emptyDir hostPath nfs persistentVolumeClaim projected secret]
hostnetwork        false     []        MustRunAs   MustRunAsRange     MustRunAs   MustRunAs   <none>     false            [configMap downwardAPI emptyDir persistentVolumeClaim projected secret]
node-exporter      false     []        RunAsAny    RunAsAny           RunAsAny    RunAsAny    <none>     false            [*]
nonroot            false     []        MustRunAs   MustRunAsNonRoot   RunAsAny    RunAsAny    <none>     false            [configMap downwardAPI emptyDir persistentVolumeClaim projected secret]
privileged         true      [*]       RunAsAny    RunAsAny           RunAsAny    RunAsAny    <none>     false            [*]
restricted         false     []        MustRunAs   MustRunAsRange     MustRunAs   RunAsAny    <none>     false            [configMap downwardAPI emptyDir persistentVolumeClaim projected secret]
```