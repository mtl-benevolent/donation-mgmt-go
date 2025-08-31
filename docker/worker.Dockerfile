FROM golang:1.24.2-bookworm AS deps

ENV NODE_MAJOR=22

RUN apt-get update -y && \
  apt-get install -y git ca-certificates make curl gnupg tzdata && \
  update-ca-certificates && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
  apt-get update -y && \
  apt-get install nodejs -y && \
  go install github.com/go-delve/delve/cmd/dlv@v1.25.2

WORKDIR /build
RUN git init --quiet

COPY go.mod go.sum package*.json Makefile ./

RUN go mod download && \
  npm ci && \
  make deps

FROM deps AS build

ARG INSTALL_BROWSERS=0

WORKDIR /build
ENV CGO_ENABLED=1

COPY ./ ./

RUN make build && \
  mkdir -p /tmp/pdfs

# Conditionally install Playwright browsers using the Makefile target
# The actual 'make install-playwright' will run the playwright CLI
RUN if [ "$INSTALL_BROWSERS" = "1" ]; then \
  make install-playwright; \
  fi

# This CMD is used for local development
CMD [ "sh", "-c", "make build && /go/bin/dlv exec ./dist/pdf-worker --headless --api-version 2 --continue --accept-multiclient --listen \"0.0.0.0:18000\"" ]

FROM debian:bookworm-slim

# --- Non-root user setup ---
ARG USERNAME=appuser
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create a non-root user and group.
# Give the user a home directory, as Playwright/Chromium might need it.
RUN groupadd --gid $USER_GID $USERNAME && \
  useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
  # Ensure /tmp is writable by the user, Chromium often uses it
  mkdir -p /tmp/pdfs && \
  chmod 1777 /tmp && \
  chmod 1777 /tmp/pdfs

COPY --from=build /go/bin/playwright /usr/local/bin/playwright

# Set DEBIAN_FRONTEND to noninteractive to avoid prompts during apt-get install
ENV DEBIAN_FRONTEND=noninteractive

ENV PLAYWRIGHT_DRIVER_PATH=/playwright/.cache/ms-playwright-go/driver
ENV PLAYWRIGHT_BROWSERS_PATH=/playwright/.cache/ms-playwright

RUN apt-get update && \
  apt-get install -y ca-certificates tzdata && \
  # Install Playwright browser and its dependencies
  # Run this as root because it installs system-level packages
  /usr/local/bin/playwright install chromium-headless-shell --with-deps && \
  apt-get clean && \
  rm -rf /var/lib/aot/lists/* && \
  # Ensure Playwright driver + browsers are executable from the non-root user
  chown -R root:$USERNAME /playwright && \
  chmod -R +rx /playwright && \
  chmod -R g+rX /playwright

WORKDIR /app

COPY --from=build /build/dist/pdf-worker ./pdf-worker

# Ensuring the binary is executable from the non-root user
RUN chown -R root:$USERNAME /app && \
  chmod +x /app/pdf-worker && \
  chmod -R g+rX /app

# --- Switch to non-root user ---
USER $USERNAME

# Setting the env values for the driver and the browsers again
ENV PLAYWRIGHT_DRIVER_PATH=/playwright/.cache/ms-playwright-go/driver
ENV PLAYWRIGHT_BROWSERS_PATH=/playwright/.cache/ms-playwright

CMD [ "/app/pdf-worker" ]
