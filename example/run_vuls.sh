#! /bin/bash
set -e

echo Using reports folder "$REPORTS_DIR"
echo Using SSH folder "$SSH_DIR"
echo Using Vulsctl folder "$VULS_DIR"

echo Running config test...
docker run --rm -it\
    -v $SSH_DIR:/root/.ssh:ro \
    -v $VULS_DIR:/vuls \
    vuls/vuls configtest \
    -log-dir=/vuls/log \
    -config=/vuls/config.toml \
    $@ 

ret=$?
if [ $ret -ne 0 ]; then
    echo Config test failed!
	exit 1
fi

echo Config test succeded, running scan...
docker run --rm -it\
    -v $SSH_DIR:/root/.ssh:ro \
    -v $VULS_DIR:/vuls \
    vuls/vuls scan \
    -log-dir=/vuls/log \
    -config=/vuls/config.toml \
    $@ 

ret=$?
if [ $ret -ne 0 ]; then
    echo Scan failed!
	exit 1
fi

echo Scan succeded, running report...
docker run --rm -it\
    -v $VULS_DIR:/vuls \
    -v $REPORTS_DIR:/reports \
    vuls/vuls report \
    -log-dir=/vuls/log \
    -format-json \
    -format-one-line-text \
    -to-localfile \
    -results-dir=/reports \
    -config=/vuls/config.toml \
    -refresh-cve \
    $@

ret=$?
if [ $ret -ne 0 ]; then
    echo Report failed!
	exit 1
fi

echo Report succeded, fixing permissions
sudo chown $USER $REPORTS_DIR -R

ret=$?
if [ $ret -ne 0 ]; then
    echo Fixing permissions failed!
	exit 1
fi

echo Vuls run succeded, now check if exporter got new metrics!
exit 0
