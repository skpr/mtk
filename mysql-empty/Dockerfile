FROM skpr/mtk-mysql:latest

# Execute a script to configure this empty MySQL instance.
ADD configure.sh /tmp/configure.sh
RUN bash /tmp/configure.sh local local local
