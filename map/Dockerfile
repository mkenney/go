FROM divan/golang:gotrace
ENV GOOS=linux
WORKDIR /go/src
RUN set -x \
    && echo 'alias ll="ls -laF"' >> /root/.bashrc \
    && echo 'alias e="exit"' >> /root/.bashrc \
    && echo 'alias cls="clear"' >> /root/.bashrc \
    # dep
    && go get -u github.com/golang/dep/cmd/dep \
    # gotrace
    && mkdir -p github.com/divan \
    && cd github.com/divan \
    && git clone https://github.com/divan/gotrace.git \
    && cd gotrace \
    && git checkout go18 \
    && dep init \
    && go install \
    # entrypoint
    && echo '#!/bin/bash' > /entrypoint.sh \
    && echo 'buildarg=$1' >> /entrypoint.sh \
    && echo 'if [ "" = "$buildarg" ]; then buildarg="."; fi' >> /entrypoint.sh \
    && echo 'echo go build -o /tmp/binary $buildarg' >> /entrypoint.sh \
    && echo 'go build -o /tmp/binary $buildarg' >> /entrypoint.sh \
    && echo '/tmp/binary 2> /tmp/trace.out' >> /entrypoint.sh \
    && echo 'if [ -s /tmp/trace.out ]; then' >> /entrypoint.sh \
    && echo '    gotrace /tmp/trace.out' >> /entrypoint.sh \
    && echo 'else' >> /entrypoint.sh \
    && echo '    echo "Trace file is empty"' >> /entrypoint.sh \
    && echo '    exit 1' >> /entrypoint.sh \
    && echo 'fi' >> /entrypoint.sh \
    && chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh
