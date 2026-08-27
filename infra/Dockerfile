FROM node:24-trixie-slim

WORKDIR /workspace/infra

# 開発ツールをインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    unzip \
    git \
    sudo \
    zsh \
    && rm -rf /var/lib/apt/lists/*

# AWS CLIをインストール
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then \
        CLI_URL="https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip"; \
    elif [ "$ARCH" = "aarch64" ]; then \
        CLI_URL="https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip"; \
    else \
        echo "Unsupported architecture: $ARCH" && exit 1; \
    fi && \
    curl "$CLI_URL" -o "awscliv2.zip" && \
    unzip -q awscliv2.zip && \
    ./aws/install && \
    rm -rf aws awscliv2.zip

# nodeユーザーでsudoをパスワードなしで使用可能にする
RUN echo "node ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/node && \
    chmod 0440 /etc/sudoers.d/node

# nodeユーザーのデフォルトシェルをzshに変更
RUN usermod -s /bin/zsh node

# zsh-autosuggestionsをインストール
RUN git clone \
    https://github.com/zsh-users/zsh-autosuggestions \
    /home/node/.zsh-autosuggestions

# zsh設定
RUN echo 'autoload -Uz vcs_info' >> /home/node/.zshrc && \
    echo 'precmd() { vcs_info }' >> /home/node/.zshrc && \
    echo "zstyle ':vcs_info:*' formats '%F{green}(%b)%f'" >> /home/node/.zshrc && \
    echo 'PROMPT="%F{cyan}%~%f ${vcs_info_msg_0_} %# "' >> /home/node/.zshrc && \
    echo 'source /home/node/.zsh-autosuggestions/zsh-autosuggestions.zsh' >> /home/node/.zshrc && \
    chown -R node:node \
        /home/node/.zshrc \
        /home/node/.zsh-autosuggestions

USER node
