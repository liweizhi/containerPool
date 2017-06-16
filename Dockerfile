FROM busybox

ENV TAG latest
ENV PATH $PATH:/go/bin:/usr/local/go/bin

COPY . /go/src/github.com/liweizhi/containerPool

WORKDIR /go/src/github.com/liweizhi/containerPool/controller

CMD ./controller -D server
