# Lightwell Experience

Lightwell Experience provides APIs and background jobs for Lightwell Network catalogs, advisories, Beacon vulnerability tracking, Lens reports, and remediation publishing via a JFrog Bridge.

## Hosting

The Lightwell Experience service (backend, UI, libraries) currently lives in the [content-services org](https://github.com/content-services). This document describes the Lightwell surface. 

Until extraction is complete, the resources defined here are shared with Content Sources.

# Architecture

Version: 1.0

Note: this document describes Lightwell Experience as it exists today. The backend service currently runs inside [content-sources-backend](https://github.com/content-services/content-sources-backend); see [Hosting](#hosting).

This application consists of these parts:

* Rest API
  * Used for our frontend or for other clients to interact with ours
* PostgreSQL Database
  * Used for storing our data (TODO: database model)
* Cron jobs
  * `import-lightwell-packages` — refreshes the package/version mirror from Pulp on a schedule
  * `sync-lightwell-advisories` — fetches OSV advisories into Postgres for the advisories API; when enabled, also emits per-org notification events and advisory-created events for JFrog Bridge
  * `sync-lightwell-vulns` — ingests vulnerabilities from Jira
* Task Worker
  * Pulls tasks from the database and executes them
* JFrog Bridge
  * Consumes advisory-created Kafka events and publishes remediations and evidence to JFrog
* Redis cache
  * Used to cache RBAC access lists, Feature Service entitled features, and repository content counts
* S3
  * Stores manifest uploads for Lens

## Deployments

When deploying within Kubernetes, we recommend:

* Minimum 3 pods for the API
  * Each pod consumes about 100 MB of memory
* Minimum 3 pods for the task worker
  * Each pod consumes about 400 MB of memory

Note: Replica and memory guidance follows the Content Sources API service and task worker deployment. Lightwell still runs inside the Content Sources processes, and those pods are sized for more than just Lightwell.

## Data Loss

* Database: If the database is lost, scheduled import/sync jobs will repopulate Lightwell catalog, advisory, and vulnerability data. Notification preferences and Lens reports cannot be rebuilt that way and must be restored from backup.
* Kafka: If queue contents are lost, pending notifications and JFrog Bridge publishes may be delayed until events are re-emitted or advisory sync runs again.
* Redis: If the cache is lost, API latency is temporarily higher until RBAC, Feature Service entitlement, and content-count caches rebuild.
* Pulp Server: If the Pulp database or content is lost, package import and content serving break until restored from backup.
* S3: If the Lens uploads bucket is lost, Lens upload objects are lost unless restored from backup.
* JFrog: Loss of catalog artifacts published by the bridge must be recovered in JFrog or by re-publishing

## External Dependencies

* Feature service
  * If unavailable, entitlement checks fail and restricted actions cannot be performed
* Jira
  * If unavailable, vulnerability sync cannot refresh data
* JFrog
  * If unavailable, bridge publishes fail
* RBAC
  * https://github.com/RedHatInsights/insights-rbac
  * If unavailable, user actions cannot be performed
* Terms service
  * If unavailable, terms routes may not be registered or checks fail. TODO: specify whether access is still granted (fail-open) or not (fail-close)
* Database
  * If unavailable, user actions and sync/import jobs cannot run
* Redis
  * If unavailable, API requests have higher latency
* Kafka broker
  * If unavailable, notifications cannot be sent and JFrog Bridge cannot consume events
* Pulp
  * If unavailable, scheduled package-mirror import, advisory syncing, content serving, and the JFrog integration are impacted; existing mirror rows remain readable until the next successful import
* S3
  * If unavailable, Lens uploads fail
* Notifications service
  * If unavailable, notifications cannot be sent
* crc-pdf-generator
  * https://github.com/RedHatInsights/pdf-generator
  * If unavailable, Beacon and Lens PDF exports fail
* Maven Central
  * If unavailable, additional details will be missing for Java packages that haven't been added to the database yet 

## Routes

Lightwell routes are exposed in two places today:

* Under Content Sources: `/api/content-sources/v1/` (OpenAPI at [/api/openapi.json](https://github.com/content-services/content-sources-backend/blob/main/api/openapi.json))
* Dedicated Lightwell API: `/api/lightwell/v0.1/` (OpenAPI at [/api/lightwell/v0.1/openapi.json](https://github.com/content-services/content-sources-backend/blob/main/api/lightwell/openapi.json))

The Content Sources OpenAPI document is available [here](https://redocly.github.io/redoc/?url=https://raw.githubusercontent.com/content-services/content-sources-backend/main/api/openapi.json).
