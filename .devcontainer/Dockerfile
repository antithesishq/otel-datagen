# Use the official VS Code base image for Debian.
# This image comes with a non-root 'vscode' user pre-configured with sudo access.
FROM mcr.microsoft.com/vscode/devcontainers/base:debian

# --- Environment Arguments ---
# Pass arguments from devcontainer.json.
ARG TZ
ARG CLAUDE_CODE_VERSION="latest"
ARG GEMINI_CLI_VERSION="latest"

# --- Environment Variables ---
# Set the timezone and prevent interactive prompts during package installation.
ENV TZ=${TZ}
ENV DEBIAN_FRONTEND=noninteractive
# Define the non-root username, matching the user in the base image.
ENV USERNAME=vscode

# --- System & Package Setup (as root) ---
USER root

# Update package lists and install essential packages.
# build-essential is needed for Homebrew.
# fish is installed here so we can set it as the default shell.
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    procps \
    curl \
    file \
    git \
    sudo \
    fish \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set fish as the default shell for the vscode user.
RUN chsh -s /usr/bin/fish ${USERNAME}

# --- User-Specific Setup (as vscode) ---
# Switch to the non-root user for subsequent installations.
USER ${USERNAME}
WORKDIR /home/${USERNAME}

# Install Homebrew as the 'vscode' user.
RUN /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Add Homebrew to the PATH. Note: The Homebrew installer script does this automatically
# for fish by creating a config file, so manual PATH setting is not strictly needed
# for the interactive shell but is good practice for scripts.
ENV PATH="/home/linuxbrew/.linuxbrew/bin:${PATH}"

RUN brew install less \
  git \
  procps \
  fzf \
  zsh
RUN brew install man-db \
  unzip
RUN brew install gh \
  gnupg2 \
  iptables \
  iproute2
RUN brew install aggregate \
  jq \
  nano \
  vim \
  bat \
  fish \
  ffmpeg \
  7zip
RUN brew install fd \
  ripgrep \
  zoxide \
  imagemagick \
  git-delta \
  mask \
  uv \
  node \
  starship
RUN brew install openjdk

# --- Node.js Global Packages ---
# Install global npm packages specified by the user.
RUN npm install -g \
    @anthropic-ai/claude-code@${CLAUDE_CODE_VERSION} \
    @google/gemini-cli@${GEMINI_CLI_VERSION}

# --- Shell & Prompt Configuration ---
# Copy the starship prompt configuration into the container.
COPY --chown=${USERNAME}:${USERNAME} starship.toml /home/${USERNAME}/.config/starship.toml

# Configure the fish shell.
# - Initialize starship prompt.
# - Add uv shell completions.
RUN mkdir -p /home/${USERNAME}/.config/fish && \
    echo 'starship init fish | source' >> /home/${USERNAME}/.config/fish/config.fish && \
    echo 'uv fish-completion >> /home/${USERNAME}/.config/fish/completions/uv.fish' >> /home/${USERNAME}/.config/fish/config.fish

# --- Final Setup ---
# Set the default working directory for the container.
WORKDIR /workspace

# Set default editor
ENV EDITOR=nano
ENV VISUAL=nano

# Set DEVCONTAINER environment variable to help with orientation
ENV DEVCONTAINER=true