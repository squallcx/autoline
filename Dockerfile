FROM registry.gitlab.com/squallcx/line
RUN apt-get update && apt-get install -y \
  unzip \
  && wget -N http://chromedriver.storage.googleapis.com/2.20/chromedriver_linux64.zip \
  && unzip chromedriver_linux64.zip \
  && chmod +x chromedriver \
  && mv -f chromedriver /usr/local/share/chromedriver \
  && ln -s /usr/local/share/chromedriver /usr/local/bin/chromedriver \
  && ln -s /usr/local/share/chromedriver /usr/bin/chromedriver 

RUN wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xf go1.5.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$PATH
RUN go version
RUN go env
RUN go get github.com/tools/godep



COPY conf/statup /root/.fluxbox/startup
COPY conf/init /root/.fluxbox/init
COPY conf/vnc_viewonly.html /opt/novnc/vnc_viewonly.html
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
COPY run.sh /run.sh

COPY code $GOPATH/src/code
