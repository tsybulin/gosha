#!/bin/sh

PRG=gosha

strip gosha

docker build -t ${PRG} .

#docker save ${PRG}>${PRG}.tar
#docker image rm ${PRG}
#docker load <${PRG}.tar
#rm ${PRG}.tar
