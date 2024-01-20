# replicator

> NOTE: Since `replicator` isn't the greatest name in the world, and [naming things](https://martinfowler.com/bliki/TwoHardThings.html) is one of the hardest things in computer science, we will delay making a decision on a name until a better one has been thought of. So for now, introducing, `replicator`....

> NOTE: Replicator is a new repository that has no real functionality yet, but the aims of what it tries to achieve are detailed below.

## What is Replicator?

Replicator is an open-source tool that aims to fill a current gap in the software supply chain world. It aims to offer engineers a more declarative way of managing their OCI artifacts, with the added capability of verification.

## Problem Context

Currently, if you were to pull an image from a remote and public registry (Dockerhub, GHCR etc), you can do a pull using your favourite CLI and start using the image instantly in whatever capacity you desire. However, say you wanted to verify the signatures and attestations of the image you've just pulled, you would normally use `cosign` directly using the CLI. This is fine for an engineer who just wants to pull an image for use locally, but as an organisation this doesn't scale.

When organisations use container OCI artefacts, there are commonly two ways that they _could_ do it:

- External pulls
  - The location of the artefact they want to use in their deployment manifests is a location that is external to the company i.e Dockerhub.
- Internal pulls
  - The location for the OCI artefact is set to an internal registry that the organisation hosts.

We won't get into the details of which is preferred between external vs internal, but for now we will focus on the internal pulls method. When an organisation pulls an artefact from a remote registry and pushes it into their own registry, there really isn't a standardised way of doing this. Some organisations have CI jobs that pull from source and push into internal registry, some have scripts run ad-hoc, either way, it's effort spent across many organisations that essentially do the same thing. Pull from one place, push into another.

Now when we start talking about attestations and provenance of artefacts, how many organisations do you think verify each artefact? In our experience we've seen no organisation actually have a handle on this as a problem without having spent copious amounts of time with an internally built solution that satisfies their needs. Policy engines allow this verification to happen on a deployment level, but in order to do the checks the signatures have to exist in the same place as the artefact, and with external registries this isn't a problem as much, but for internal registries, how do you get those signatures.. replicated?

This is where replicator offers a hand.

## So, how does replicator help solve the problem?

Replicator is an application that will be configured to run through all declared artefacts that an organisation wants to use, and aims to do the following:

- Pulls artefact
- Verifies artefact (signatures, attestations, provenance)
- Pushes artefact to internal registry with its provenance and attestation information

It aims to be simple in nature but understands there will be complexity in the detail. However, using Replicator shall allow organisations to declare their OCI artefact estate in Git whilst it goes through and ensures that the artefacts declared are "replicated" into the internal registry ensuring the verification checks are done.
