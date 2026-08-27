# -----------------------------------
# 共通設定 (Target: base)
# -----------------------------------

FROM node:24-trixie-slim AS base

# UID/GIDの受け取り
ARG UID=1000
ARG GID=1000

# パッケージリストの更新とツールのインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    zsh \
    zip \
    && rm -rf /var/lib/apt/lists/*

# nodeユーザーのUID/GIDを変更
RUN usermod -u ${UID} node && \
    groupmod -g ${GID} node

# nodeユーザーのデフォルトシェルをzshに変更
RUN usermod -s /bin/zsh node

# ディレクトリ
WORKDIR /workspace/frontend

# zshにGitブランチ表示と色付けを設定
RUN echo 'autoload -Uz vcs_info' >> /home/node/.zshrc && \
    echo 'precmd() { vcs_info }' >> /home/node/.zshrc && \
    echo "zstyle ':vcs_info:*' formats '%F{green}(%b)%f'" >> /home/node/.zshrc && \
    echo 'PROMPT="%F{cyan}%~%f ${vcs_info_msg_0_} %# "' >> /home/node/.zshrc

# zsh-autosuggestionsをインストール
RUN git clone \
    https://github.com/zsh-users/zsh-autosuggestions \
    /home/node/.zsh-autosuggestions && \
    echo 'source /home/node/.zsh-autosuggestions/zsh-autosuggestions.zsh' \
    >> /home/node/.zshrc

# 所有者をnodeユーザーに変更
RUN chown -R node:node \
    /home/node/.zshrc \
    /home/node/.zsh-autosuggestions \
    /workspace/frontend

# -----------------------------------------
# 開発用 (Target: dev)
# -----------------------------------------

FROM base AS dev

# rootを明示
USER root

# ブラウザをインストールする場所を固定
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

# デバッグ用ツール・依存関係をインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    iputils-ping \
    net-tools \
    procps \
    lsof \
    libnss3 \
    libnspr4 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxfixes3 \
    libxrandr2 \
    libgbm1 \
    libasound2t64 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libxcursor1 \
    libgtk-3-0 \
    libgdk-pixbuf-xlib-2.0-0 \
    libcairo-gobject2 \
    libxss1 \
    && rm -rf /var/lib/apt/lists/*

# Playwright用ディレクトリを作成
RUN mkdir -p /ms-playwright && \
    chown -R node:node /ms-playwright

# 設定ファイルをコピー
COPY --chown=node:node package*.json ./

# ユーザー切り替え
USER node

# npm依存関係をインストール
RUN npm ci

# Chromiumブラウザをインストール
RUN npx playwright install chromium

# Git所有権エラー回避
RUN git config --global --add safe.directory '*'

# ソースコードをコピー
COPY --chown=node:node . .

# Vite開発サーバーのポートを開放
EXPOSE 5173

# -----------------------------------------
# ビルド＆ZIP用 (Target: builder)
# -----------------------------------------

FROM base AS builder

USER node

COPY --chown=node:node package*.json ./
RUN npm ci

COPY --chown=node:node . .

# Viteでビルドしてdistフォルダを作成
RUN npm run build

# ビルド成果物（dist）をZIPに固める
RUN cd dist && zip -r ../frontend-dist.zip .
