# ------------------------------------
# 開発用 (Target: dev)
# ------------------------------------
FROM golang:1.26.3-trixie AS dev

# UID/GIDの受け取り
ARG UID=1000
ARG GID=1000

# 依存関係・AWS CLI・デバッグツールをインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    build-essential \
    zsh \
    curl \
    unzip \
    less \
    groff \
    jq \
    iputils-ping \
    net-tools \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf aws awscliv2.zip \
    && rm -rf /var/lib/apt/lists/*

# ユーザー作成
RUN groupadd -g ${GID} appgroup && \
    useradd -u ${UID} -g ${GID} -s /bin/zsh -m appuser

# zshにGitブランチ表示と色付けを設定
RUN echo 'autoload -Uz vcs_info' >> /home/appuser/.zshrc && \
    echo 'precmd() { vcs_info }' >> /home/appuser/.zshrc && \
    echo "zstyle ':vcs_info:*' formats '%F{green}(%b)%f'" >> /home/appuser/.zshrc && \
    echo 'PROMPT="%F{cyan}%~%f ${vcs_info_msg_0_} %# "' >> /home/appuser/.zshrc && \
    chown appuser:appgroup /home/appuser/.zshrc

# ディレクトリ
WORKDIR /workspace/backend

# モジュールキャッシュの権限設定
RUN mkdir -p /go/pkg/mod /go/bin && \
    chown -R appuser:appgroup /workspace/backend /go/pkg /go/bin

# ユーザー切り替え
USER appuser

# パス設定
ENV PATH=$PATH:/go/bin

# Air + golangci-lint インストール
RUN --mount=type=cache,target=/go/pkg/mod,uid=${UID},gid=${GID} \
    go install github.com/air-verse/air@v1.65.3 && \
    go install golang.org/x/tools/cmd/goimports@v0.48.0 && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

# モジュールダウンロード
COPY --chown=appuser:appgroup go.mod go.sum* ./
RUN go mod download

# ソースコードをコピー
COPY --chown=appuser:appgroup . .

CMD ["air", "-c", ".air.toml"]

# ------------------------------------
# ビルド用 (Target: builder)
# ------------------------------------
FROM golang:1.26.3-trixie AS builder
WORKDIR /workspace/backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# バイナリとしてビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api/main.go

# ------------------------------------
# 本番用 (Target: runner)
# ------------------------------------
FROM debian:13-slim AS runner
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/server .
USER 1000:1000
CMD ["./server"]
