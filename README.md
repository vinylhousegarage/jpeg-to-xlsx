# jpeg-to-xlsx

## 1. 概要

- スマートフォンで撮影した画像から文字情報を抽出し、Excelファイルへ変換するWebアプリです。
- ExcelファイルをダウンロードするURLを、SlackのDMへ通知します。

## 2. デモ（YouTubeショート / 音声あり）

[![jpeg-to-xlsxのデモ動画](https://img.youtube.com/vi/XREA43I91Rg/maxresdefault.jpg)](https://www.youtube.com/shorts/XREA43I91Rg)

**※公開時に撮影したデモです。現在、Webアプリの公開は停止しています。**

## 3. システム構成

![システム構成図](docs/system-diagram.png)

## 4. ログイン経路

```mermaid
sequenceDiagram
  autonumber
  participant User as ユーザー
  participant Browser as スマートフォン<br/>（ブラウザ）
  participant CloudFront
  participant S3Web as S3<br/>（配信用）
  participant API as API Gateway<br/>（HTTP API）
  participant BFF as Lambda<br/>（API実行）
  participant StateDB as DynamoDB<br/>（Cognito state保存）
  participant Cognito
  participant Google as Google OAuth
  participant SessionDB as DynamoDB<br/>（セッション保存）

  User->>Browser: jpeg-to-xlsxへアクセス
  Browser->>CloudFront: Webアプリを要求
  CloudFront->>S3Web: 静的ファイルを取得
  S3Web-->>CloudFront: HTML・JS・CSS
  CloudFront-->>Browser: Reactアプリを返す

  Browser->>CloudFront: /api/auth/login
  CloudFront->>API: /api/*を転送
  API->>BFF: ログイン開始
  BFF->>StateDB: state・nonce・PKCE verifierを保存
  BFF-->>Browser: Cognitoへリダイレクト

  Browser->>Cognito: 認証要求
  Cognito->>Google: Google認証へ転送
  User->>Google: アカウント選択・本人確認
  Google-->>Cognito: 認証成功
  Cognito-->>Browser: code・state付きでLamda（API実行）へ戻す

  Browser->>CloudFront: /api/auth/callback
  CloudFront->>API: /api/*を転送
  API->>BFF: code・state
  BFF->>StateDB: stateで一時情報を取得
  StateDB-->>BFF: nonce・PKCE verifier
  BFF->>Cognito: code・verifierでトークン交換
  Cognito-->>BFF: IDトークン
  BFF->>BFF: nonce・IDトークンを検証
  BFF->>SessionDB: セッションを保存
  BFF-->>Browser: Session Cookieを設定
```

## 5. アップロード用S3署名付きURL取得経路

```mermaid
sequenceDiagram
  autonumber
  participant User as ユーザー
  participant Browser as スマートフォン<br/>（ブラウザ）
  participant CloudFront
  participant API as API Gateway<br/>（HTTP API）
  participant BFF as Lambda<br/>（API実行）
  participant SessionDB as DynamoDB<br/>（セッション保存）
  participant S3Input as S3<br/>（アップロード用）

  User->>Browser: 画像を撮影・送信
  Browser->>CloudFront: S3署名付きURLを要求<br/>Session Cookie付き
  CloudFront->>API: /api/*を転送
  API->>BFF: URL発行要求
  BFF->>BFF: CookieのSession IDをハッシュ化
  BFF->>SessionDB: id_hashでセッションを取得
  SessionDB-->>BFF: cognito_sub・有効期限
  BFF->>BFF: ログイン状態と入力値を確認
  BFF->>S3Input: アップロード用S3署名付きURLを生成
  S3Input-->>BFF: S3署名付きURL
  BFF-->>API: URL・オブジェクトキー
  API-->>CloudFront: APIレスポンス
  CloudFront-->>Browser: S3署名付きURLを返す
```

## 6. Slack通知経路

```mermaid
sequenceDiagram
  autonumber
  participant User as ユーザー
  participant Browser as スマートフォン<br/>（ブラウザ）
  participant S3Input as S3<br/>（アップロード用）
  participant Processor as Lambda<br/>（イベント駆動）
  participant Bedrock as Amazon Bedrock<br/>（Claude Sonnet 4.6）
  participant S3Output as S3<br/>（ダウンロード用）
  participant TokenDB as DynamoDB<br/>（Slackトークン）
  participant Slack as Slack API

  Browser->>S3Input: S3署名付きURLで画像をアップロード
  S3Input-->>Browser: アップロード成功

  S3Input-->>Processor: S3イベントで非同期起動
  Processor->>S3Input: 画像を取得
  S3Input-->>Processor: JPEGデータ

  Processor->>Bedrock: 画像解析を要求
  Bedrock-->>Processor: JSON形式の抽出結果
  Processor->>Processor: JSONをXLSXへ変換

  Processor->>S3Output: XLSXファイルを保存
  S3Output-->>Processor: 保存完了
  Processor->>Processor: ダウンロード用S3署名付きURLを生成

  Processor->>TokenDB: Slackトークンを取得
  TokenDB-->>Processor: Skackトークン
  Processor->>Slack: XLSXのS3署名付きURLをDM送信
  Slack-->>User: ダウンロードリンクを通知
```

## 7. DynamoDBテーブル構成

```mermaid
erDiagram
  COGNITO_OAUTH_STATES {
    string state PK
    string code_verifier
    string nonce
    number created_at
    number expires_at
  }

  AUTH_SESSIONS {
    string id_hash PK
    string cognito_sub
    number created_at
    number expires_at
  }

  SLACK_TOKENS {
    string cognito_sub PK
    string team_id
    string access_token
    string bot_user_id
    string user_id
    string channel_id
    string updated_at
  }
```

## 8. 開発目的

- 紙媒体に記載された文字情報をExcelに変換することで、データ入力や転記に伴う作業の効率化を支援することを目的としています。

## 9. 技術スタック

### 開発基盤
  | 項目 | 技術 |
  | :--- | :--- |
  | 開発環境 | Docker |
  | OS | Debian 13 |
  | ソース管理 | Git |
  | リポジトリ | GitHub |
  | CI/CD | GitHub Actions |

### バックエンド
  | 項目 | 技術 |
  | :--- | :--- |
  | 開発言語 | Go 1.26.3 |
  | 認可連携 | Slack OAuth |
  | 通知連携 | Slack Web API |
  | 生成AI基盤 | Amazon Bedrock |
  | 生成AIモデル | Claude Sonnet 4.6 |

### フロントエンド
  | 項目 | 技術 |
  | :--- | :--- |
  | 開発言語 | TypeScript 6.0.3 |
  | UIライブラリ | React 19.2.7 |
  | ビルドツール | Vite 8.1.0 |
  | ビルド環境 | Node.js 24.19.0 |

### インフラ（AWS）
  | 項目 | 技術 |
  | :--- | :--- |
  | 開発言語 | TypeScript 6.0.3 |
  | IaC | CloudFormation |
  | IaCフレームワーク | AWS CDK (aws-cdk-lib 2.263.0) |
  | IaC CLI | AWS CDK CLI 2.1127.0 |
  | AWS CDK実行環境 | Node.js 24.19.0 |
  | 配信基盤 | CloudFront |
  | 認証基盤 | Cognito |
  | オブジェクトストレージ | S3 |
  | API基盤 | API Gateway |
  | APIタイプ | HTTP API |
  | 実行基盤 | Lambda |
  | データベース | DynamoDB |
  | 秘匿情報管理 | AWS Secrets Manager |

## 10. 技術選定

- AWSが提供するIaCサービス、CloudFormationを採用しました。

- CloudFormationテンプレートを生成することができるIaCフレームワーク、AWS CDKを採用しました。

- AWS公式ドキュメントやAWS CDKのTypeScript向けサンプルが充実しているため、インフラの開発言語にTypeScriptを採用しました。

- スマートフォンでの撮影・プレビュー・アップロード・結果表示のUIと画面状態を管理するため、UIライブラリであるReactを採用しました。

- コンパイル型言語による実行性能とLambdaとの親和性を考慮し、バックエンドの開発言語にGoを採用しました。

- Googleアカウントによるログインを実現し、アプリケーションの認証情報を一元管理するため、認証基盤としてAmazon Cognitoを採用しました。

- 画像内の文字情報をJSON形式に構造化する精度と運用コストのバランスを考慮し、生成AIモデルにClaude Sonnet 4.6を採用しました。

- CloudFrontを採用し、フロントエンドの配信とバックエンドへのアクセスを同一ドメインにまとめました。

- S3を採用し、アップロード用S3署名付きURLを用いてブラウザから撮影画像を直接アップロードする構成としました。

- Slackを採用し、画像解析完了後にダウンロード用S3署名付きURLをDMへ通知する構成としました。

- Slackとの認可連携を行うため、Slack OAuthを採用しました。

- Slack OAuthのクライアントシークレットを安全に管理するため、AWS Secrets Managerを採用しました。

- Slack OAuthで取得したアクセストークンをサーバーレスで管理するため、データベースにDynamoDBを採用しました。

## 11. ライセンス

- 本リポジトリは [MIT License](./LICENSE) のもとで公開しています。
