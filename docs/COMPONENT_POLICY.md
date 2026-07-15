# Component Policy

This document defines the criteria for adding new OpenTelemetry Collector components (receivers, processors, exporters, extensions, providers) to sacloud-otel-collector.

(日本語訳は[下記](#コンポーネントポリシー日本語訳)にあります)

## Purpose

The purpose of this collector is to collect telemetry on SAKURA Cloud environments and forward it to [SAKURA Cloud Monitoring Suite](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/index.html) and other common destinations. This distribution does not aim to reproduce the general-purpose `otelcol-contrib` distribution.

## Criteria

### 1. Relevance to the purpose (required)

The component must fall into one of the following categories:

- **Platform / infrastructure telemetry** — collecting telemetry from SAKURA Cloud platform features and APIs, or from the OS and container runtime layer (not tied to a specific application).
- **Generic ingestion** — receiving telemetry from applications running on customer VMs via standard protocols (OTLP, Prometheus, etc.). Application-specific receivers are out of scope; applications should emit telemetry through these generic protocols.
- **Forwarding** — exporting telemetry to SAKURA Cloud Monitoring Suite or other widely used observability backends.

### 2. Generality (required)

- Receivers dedicated to a specific middleware or application (e.g. nginx, mysql, redis) are not added in principle. If telemetry can be collected with existing components such as the Prometheus, filelog, or OTLP
  receivers, use those instead.
- Do not add a new processor for transformations that can be achieved with the configuration of existing components (e.g. the transform processor).

### 3. Stability

- The component must be beta or higher for the target signal(s) in principle. Check the stability table in the component's README or `metadata.yaml`.
- Alpha components are acceptable only when there is no alternative, a concrete use case exists, and the component is continuously maintained (precedents: journald, windowseventlog, deltatocumulative).

### 4. Source of the component

Preference order: core > contrib > third-party.

Third-party components (precedent: mackerelotlpexporter) are acceptable only when the service has a relationship with the SAKURA ecosystem and the component is actively maintained to follow OpenTelemetry Collector releases.

### 5. Platform support

The component must build on all release targets (Linux, macOS, Windows). Components that only work on a specific OS are acceptable as long as they build on the other platforms (precedents: windowseventlogreceiver, journaldreceiver).

### 6. Trigger for addition

Add a component only after a concrete use case is presented, e.g. via an issue. Do not add components speculatively because they "might be useful".

## How to add a component

See the [Contributing](../README.md#contributing) section of the README.

---

# コンポーネントポリシー(日本語訳)

このドキュメントは、sacloud-otel-collector に新しい OpenTelemetry Collector コンポーネント(receiver, processor, exporter, extension, provider)を追加する際の基準を定めます。

## 目的

この collector の目的は、さくらのクラウド環境でテレメトリを収集し、
[SAKURA Cloud Monitoring Suite](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/index.html)
やその他の一般的な送信先へ転送することです。汎用ディストリビューション(`otelcol-contrib`)の再現は目指しません。

## 基準

### 1. 目的適合性(必須)

コンポーネントが以下のいずれかに該当すること。

- **プラットフォーム・インフラのテレメトリ** — さくらのクラウドのプラットフォーム機能・API、または OS やコンテナランタイムなど特定アプリケーションに依存しない基盤レイヤーからテレメトリを収集するもの。
- **汎用プロトコルによる受信** — 顧客 VM 上のアプリケーションから OTLP や Prometheus などの標準プロトコルでテレメトリを受け取るもの。アプリケーション専用の receiver は対象外とし、アプリケーション側でこれらの汎用プロトコルを通じてテレメトリを送信する想定とする。
- **転送** — SAKURA Cloud Monitoring Suite やその他の広く使われている観測バックエンドへテレメトリを送信するもの。

### 2. 汎用性(必須)

- 特定のミドルウェア・アプリケーション専用の receiver(例: nginx, mysql, redis)は原則追加しません。Prometheus, filelog, OTLP receiver など既存のコンポーネントで収集できる場合はそちらを使用します。
- 既存コンポーネントの設定で実現可能な処理(例: transform processor で書けるもの)のために新しい processor を追加しません。

### 3. 安定性

- 対象シグナルについて原則 beta 以上 であること。コンポーネントの README または `metadata.yaml` の stability 表を確認してください。
- alpha は「代替手段がなく、具体的なユースケースがある」「メンテナンスが継続的に行われている」場合のみ許容します(前例: journald, windowseventlog, deltatocumulative)。

### 4. 供給元

優先順位: core > contrib > サードパーティ。

サードパーティ製(前例: mackerelotlpexporter)は、さくらエコシステムと連携関係があり、OpenTelemetry Collector 本体のリリースに追従してメンテナンスされている場合のみ許容します。

### 5. プラットフォーム

リリース対象(Linux / macOS / Windows)すべてでビルドが通ること。特定 OS でしか動作しないコンポーネントでも、他の OS でビルド可能であれば許容します(前例: windowseventlogreceiver, journaldreceiver)。

### 6. 追加の起点

Issue 等で具体的なユースケースが示されてから追加します。「あると便利かもしれない」という推測での先回り追加はしません。

## コンポーネントの追加手順

README の [Contributing](../README.md#contributing) セクションを参照してください。
