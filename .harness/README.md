# Run tests on Harness CI

This application can be built and tested with [Harness CI](https://www.harness.io/products/continuous-integration).

This guide covers how to use [AIDA](https://www.harness.io/products/aida) and [Remote Debug](https://developer.harness.io/docs/continuous-integration/troubleshoot-ci/debug-mode/) with Harness CI.

## Setting up this pipeline on Harness CI Hosted Builds

1. Create a [GitHub Account](https://github.com) or use an existing account

2. Fork [this repository](https://github.com/harness-community/go-fruits-api/fork) into your GitHub account

3.
    a. If you are new to Harness CI, signup for [Harness CI](https://app.harness.io/auth/#/signup)
      * Select the `Continuous Integration` module and choose the `Starter pipeline` wizard to create your first pipeline using the forked repo from #2.

    b. If you are an existing Harness CI user, create a new pipeline to use the cloud option for infrastructure and setup the PR trigger.

4. To enable Remote Debug in your Harness account, check the [Remote Debug](https://developer.harness.io/docs/continuous-integration/troubleshoot-ci/debug-mode/) documentation.

5. To enable AIDA in your Harness account, see [Get started with AIDA for CI](https://developer.harness.io/docs/continuous-integration/troubleshoot-ci/aida/#get-started-with-aida-for-ci).

6. Insert this YAML into your pipeline's `stages` section.

```yaml
    - stage:
        name: test
        identifier: test
        description: ""
        type: CI
        spec:
          cloneCodebase: true
          platform:
            os: Linux
            arch: Amd64
          runtime:
            type: Cloud
            spec: {}
          execution:
            steps:
              - step:
                  type: Background
                  name: postgres
                  identifier: postgres
                  spec:
                    image: postgres:14.4-alpine
                    shell: Sh
                    envVariables:
                      POSTGRES_PASSWORD: password
                      POSTGRES_USER: demo
                      POSTGRES_DB: demodb
                    portBindings:
                      "5432": "5432"
              - step:
                  type: Run
                  name: test
                  identifier: test
                  spec:
                    shell: Sh
                    command: |-
                      go test -race ./... -v
                    envVariables:
                      POSTGRES_HOST: postgres
                      POSTGRES_PORT: "5432"
                      POSTGRES_USER: demo
                      POSTGRES_PASSWORD: password
                      POSTGRES_DB: demodb
                      FRUITS_DB_TYPE: pgsql
```

7. Save your changes, then click "Run" to run the pipeline.

    The pipeline will fail with this error:

    ```
    dial tcp: lookup postgres: no such host
    ```

8. Click the __Ask AIDA__ button beneath the log output. AIDA will examine the log output, then click the __View__ button to see suggested fixes.

9. Open the failed pipeline exution in a different browser tab, select __More Options__ (⋮), then select __Re-run in Debug Mode__.

10. Create a [Harness API key](https://developer.harness.io/docs/platform/automation/api/add-and-manage-api-keys) with pipeline execution permissions. You must have `pipeline execution` permissions at the [account scope](https://developer.harness.io/docs/platform/role-based-access-control/rbac-in-harness#permissions-hierarchy-scopes) in order for the token to have those permissions.

11. The pipeline will pause at the failing `test` step and show the following message:

    ```
    Connecting to ssh.harness.io...
    In order to ssh into your env please copy the following command:
    ssh  {harness pat}:<your-harness-account-ID>:<random-session-token>@<subdomain>.harness.io 
    where {harness_pat} should be replaced with a valid harness access token with execute permissions
    ```

    Copy the `ssh` line into your terminal, replace `{harness pat}` with the API key generated in step 10 above.

    You will now have an SSH session into your pipeline step, the prompt will look something like this:

    ```
    root@6e7e65fb23628de1:/harness# 
    ```

12. Reproduce the `test` pipeline step failure by running the test command in your SSH session:

    ```
    go test -race ./... -v
    ```

    You will see the same `dial tcp: lookup postgres: no such host` error.

13. See the __Root cause__ and __Remediation__ information that AIDA generated in step 8. It should say something like this:

    > To fix the issue, ensure that the hostname "postgres" is resolvable from the machine where the tests are being run.

    The failing `test` pipeline sets the environment variable `POSTGRES_HOST: postgres`, which is causing the failure.

    Environment variables set for the `test` step are also set in the debug session. Try echoing variables like `POSTGRES_HOST` and `POSTGRES_USER`.

    For example, `echo $POSTGRES_HOST` will print `postgres`, `echo $POSTGRES_USER` will print `demo`.

15. The [Background step settings](https://developer.harness.io/docs/continuous-integration/use-ci/manage-dependencies/background-step-settings/) documentation describes how to communicate with background steps.

    Since `postgres` is background step running in a Docker container, and the `test` step is running on the host, `POSTGRES_HOST` must be `localhost`, not `postgres`.

16. Test the fix in the debug session.

    In your debug session, set `POSTGRES_HOST` to `localhost` and run the test command again.

    ```
    export POSTGRES_HOST=localhost
    go test -race ./... -v
    ```

    This time the tests will complete successfully.

    When you have finished debugging, select __Abort Pipeline__ from the __More Options__ (⋮) menu.

17. Apply the fix to your pipeline.

    Make this change to your pipeline's `test` step.

    ```diff
                     envVariables:
    -                  POSTGRES_HOST: postgres
    +                  POSTGRES_HOST: localhost
    ```

    | ℹ️ Note |
    |---------|
    | An alternate solution would be to change the `test` step to run in a Docker container, then `POSTGRES_HOST` would not need to be changed. When both steps use Docker containers, steps can communicate with the background step via [Docker networking](https://docs.docker.com/network/), which uses the background step identifier (`postgres` in this case). |

18. Save the change and run the pipeline.

    This time, the pipeline will complete successfully.
