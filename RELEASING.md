# Release Process

This project uses the [`multimod` releaser
tool](https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod)
to manage releases. This document will walk you through how to perform a
release using this tool for this repository.

## Start a release

1. Find out which modules were changed since the last release.
2. Decide the new version of each of those modules, following semantic versioning
3. Look at the dependency graph from the [README](./README.md) and also bump the bugfix version part of all the modules that depend on the previously identified modules

For example, if you created a new feature in `x/events` and in `logger`, you would increase the
minor version of those 2 modules and also increase the bugfix versions of
`otel`, `camunda` and `middleware` as those 3 are depending on the other 2.

### Create a release branch

Update the versions of the module sets you have identified in `versions.yaml`.
Commit this change to a new release branch.

### Update module set versions

Set the version for all the module sets you have identified to be released.

```sh
make prerelease MODSET=<module set>
```

where `<module set>` is something like `logger` or `xevents` for example.

This will use `multimod` to upgrade the module's versions and create a new
"prerelease" branch for the changes. Verify the changes that were made.

```sh
git diff HEAD..prerelease_<module set>_<version>
```

Fix any issues if they exist in that prerelease branch, and when ready, merge
it into your release branch.

```sh
git merge prerelease_<module set>_<version>
```

### Make a Pull Request

Push your release branch and create a pull request for the changes. Be sure to
include the curated changes you included in the changelog in the description.
Especially include the change PR references, as this will help show viewers of
the repository looking at these PRs that they are included in the release.

## Tag a release

Once the Pull Request with all the version changes has been approved and merged
it is time to tag the merged commit.

***IMPORTANT***: It is critical you use the same tag that you used in the
Pre-Release step! Failure to do so will leave things in a broken state. As long
as you do not change `versions.yaml` between pre-release and this step, things
should be fine.

1. For each module set that will be released, run the `add-tags` make target
   using the `<commit-hash>` of the commit on the main branch for the merged
   Pull Request.

   ```sh
   make add-tags MODSET=<module set> COMMIT=<commit hash>
   ```

   It should only be necessary to provide an explicit `COMMIT` value if the
   current `HEAD` of your working directory is not the correct commit.

   You have to have gpg installed and a GPG key set up and configured in GitHub:
   https://docs.github.com/en/authentication/managing-commit-signature-verification

   1. Generate a GPG key with an email address that matches with the email address of github and
      you use when committing
   2. Export the GPG key and configure it in Github
   3. Configure your git command line to use the GPG key
   4. Set this environment variable so git is able to show the GPG passphrase prompt
      ```sh
      export GPG_TTY=$(tty)
      ```

3. If the version of Camunda was changed, you have to create another tag manually:

   ```sh
   git tag camunda/v2.0.5
   ```
   assuming that `v2.0.5` is the version that was just created 

4. Push tags to the upstream remote. Make sure you
   push all sub-modules as well.

   ```sh
   git push origin --tags
   ```

## Release

Finally create a Release on GitHub. If you are releasing multiple versions for
different module sets, be sure to use the stable release tag but be sure to
include each version in the release title (i.e. `Release v1.0.0/v0.25.0`). The
release body should include all the curated changes for this release.
